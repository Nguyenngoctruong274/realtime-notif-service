package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
	logger_custom "yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/utils/constants"
	"yes4all/ads-noti-api/pkg/websocket_service"
	"yes4all/ads-noti-api/services/ads-noti/model/entity"
	"yes4all/ads-noti-api/services/ads-noti/model/enum/enum_notification"
	"yes4all/ads-noti-api/services/ads-noti/model/request"
	"yes4all/ads-noti-api/services/ads-noti/repository"

	"github.com/gorilla/websocket"
)

type WebsocketUsecase interface {
	HandleConnection(ctx context.Context, conn *websocket.Conn, username string) (err error)
	AddNotification(ctx context.Context, req request.MessageRequest) (err error)
	SubData(ctx context.Context, req request.MarkAsReadRequest) (err error)
}

type websocketUsecase struct {
	websocketService *websocket_service.Server
	notificationRepo repository.NotificationRepo
}

func NewWebsocketUsecase(
	notificationRepo repository.NotificationRepo,
) WebsocketUsecase {
	return &websocketUsecase{
		websocketService: websocket_service.GetServerSocket(),
		notificationRepo: notificationRepo,
	}
}

func (w *websocketUsecase) HandleConnection(ctx context.Context, conn *websocket.Conn, username string) (err error) {
	entry := logger_custom.NewLogger().WithKeyword(ctx, "HandleConnection")

	defer conn.Close()
	//validate
	if len(username) == 0 {
		err = errors.New("username is empty")
		w.websocketService.SendError(conn, err.Error())
		return
	}
	w.websocketService.Mu.Lock()
	if !w.websocketService.IsConnect[username] {
		notificationDB, mErr := w.notificationRepo.GetListByUser(ctx, username)
		if mErr != nil {
			err = mErr
			w.websocketService.SendError(conn, err.Error())
			return
		}

		temp := []request.Notification{}
		for _, item := range notificationDB {
			isLink := false
			if item.Notification.Type == "TICKET" && strings.Contains(item.Content, "has assigned") {
				isLink = true
			}

			temp = append(temp, request.Notification{
				ID:             item.ID,
				Priority:       item.Priority,
				IsMarked:       item.IsMarked,
				Content:        item.Content,
				ImageUrl:       item.ImageUrl,
				ProfileID:      item.Notification.ProfileID,
				ObjectID:       item.Notification.ObjectID,
				Type:           item.Notification.Type,
				NotificationID: item.Notification.ID,
				CreatedAt:      *item.CreatedAt,
				FromUser:       item.EmailBy,
				IsLink:         isLink,
				Description:    item.Description,
			})
		}
		w.websocketService.Notifications[username] = append(w.websocketService.Notifications[username], temp...)
	}
	w.websocketService.Mu.Unlock()

	// get Notifications
	w.websocketService.Subscribe(conn, username)
	// Start goroutines to read and write messages.

	// create channel to signal client health
	done := make(chan struct{})
	go w.writePump(conn, username, done)
	w.readPump(conn, username, done)
	entry.Info()
	return nil
}

func (w *websocketUsecase) AddNotification(ctx context.Context, req request.MessageRequest) (err error) {
	entry := logger_custom.NewLogger().WithKeyword(ctx, "AddNotification")
	//validate
	username := req.Username
	if len(username) == 0 {
		err = errors.New("username is empty")
		entry.WithError(err).Error()
		return
	}
	if len(req.Data) == 0 {
		err = errors.New("data is empty")
		entry.WithError(err).Error()
		return
	}
	//save
	req.Action = enum_notification.AddAction().Data
	//
	for i, item := range req.Data {
		if item.Type == "TICKET" && strings.Contains(item.Content, "has assigned") {
			req.Data[i].IsLink = true
		}
	}

	//add data
	w.websocketService.Mu.Lock()
	unRead := 0
	if w.websocketService.IsConnect[username] {
		w.websocketService.Notifications[username] = append(
			req.Data, w.websocketService.Notifications[username]...)
		notifications, ok := w.websocketService.Notifications[username]
		if ok {
			for _, item := range notifications {
				if !item.IsMarked {
					unRead++
				}
			}
		}
	}
	w.websocketService.Mu.Unlock()

	req.Unread = &unRead
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	//sent message
	w.websocketService.Publish(username, reqBytes)

	entry.Info()
	return
}

// // readPump process incoming messages and set the settings
func (w *websocketUsecase) readPump(conn *websocket.Conn, username string, done chan<- struct{}) {
	// set limit, deadline to read & pong handler
	conn.SetReadLimit(websocket_service.MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(websocket_service.PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(websocket_service.PongWait))
		return nil
	})

	// message handling
	for {
		// read incoming message
		_, msg, err := conn.ReadMessage()
		// if error occured
		if err != nil {
			// remove from the client
			w.websocketService.Unsubscribe(conn, username)
			// set health status to unhealthy by closing channel
			close(done)
			// stop process
			break
		}
		// if no error, process incoming message
		w.ProcessMessage(conn, username, msg)
	}
}

// writePump sends ping to the client
func (w *websocketUsecase) writePump(conn *websocket.Conn, username string, done <-chan struct{}) {
	// create ping ticker
	ticker := time.NewTicker(websocket_service.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// send ping message
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(websocket_service.WriteWait))
			if err != nil {
				// if error sending ping, remove this client from the server
				w.websocketService.Unsubscribe(conn, username)
				// stop sending ping
				return
			}
		case <-done:
			// if process is done, stop sending ping
			return
		}
	}
}

func (w *websocketUsecase) SubData(ctx context.Context, req request.MarkAsReadRequest) (err error) {
	entry := logger_custom.NewLogger().WithKeyword(ctx, "NotificationHandler.writePump")
	now := time.Now()
	emailBy := req.Data.FromUser
	emailTo := req.Username

	if req.Username == "" {
		err = errors.New("user not exist. Please check again")
		entry.WithError(err).Error()
		return
	}

	schema := &entity.Notification{
		ProfileID: req.Data.ProfileID,
		ObjectID:  req.Data.ObjectID,
		Type:      req.Data.Type,
		Base: entity.Base{
			CreatedAt: &now,
			UpdatedAt: &now,
			CreatedBy: &emailBy,
			UpdatedBy: &emailBy,
		},
		UserNotifications: []entity.UserNotification{
			{
				EmailBy:     emailBy,
				EmailTo:     emailTo,
				Priority:    req.Data.Priority,
				IsMarked:    req.Data.IsMarked,
				Content:     req.Data.Content,
				ImageUrl:    req.Data.ImageUrl,
				Description: req.Data.Description,
				Base: entity.Base{
					CreatedAt: &now,
					UpdatedAt: &now,
					CreatedBy: &emailBy,
					UpdatedBy: &emailBy,
				},
			},
		},
	}
	err = w.notificationRepo.Create(ctx, schema)
	if err != nil {
		entry.WithError(err).Error()
		return
	}

	tempData := request.MessageRequest{}
	tempData.Username = req.Username
	tempData.Action = enum_notification.AddAction().Data
	for _, item := range schema.UserNotifications {
		tempData.Data = append(tempData.Data, request.Notification{
			ID:          item.ID,
			Priority:    item.Priority,
			IsMarked:    item.IsMarked,
			Content:     item.Content,
			FromUser:    item.EmailBy,
			CreatedAt:   *item.CreatedAt,
			ObjectID:    schema.ObjectID,
			Type:        schema.Type,
			ProfileID:   schema.ProfileID,
			ImageUrl:    item.ImageUrl,
			Description: item.Description,
		})
	}

	w.AddNotification(ctx, tempData)
	return
}

// ProcessMessage handle message according to the action type
func (w *websocketUsecase) ProcessMessage(conn *websocket.Conn, username string, msg []byte) {
	// parse message
	m := request.MarkAsReadRequest{}
	if err := json.Unmarshal(msg, &m); err != nil {
		err = errors.New("server: Invalid msg")
		w.websocketService.SendError(conn, err.Error())
		return
	}

	// convert all action to lowercase and remove whitespace
	action := strings.TrimSpace(strings.ToLower(m.Action))
	switch action {
	case constants.Update:
		w.UpdateNotification(conn, username, m)
	default:
		err := errors.New("server: Action unrecognized")
		w.websocketService.SendError(conn, err.Error())
	}

}

// UpdateNotification update mark a read for notification
func (s *websocketUsecase) UpdateNotification(conn *websocket.Conn, username string, req request.MarkAsReadRequest) {
	//validate
	if err := s.validate(conn, req, username); err != nil {
		s.websocketService.SendError(conn, err.Error())
		logger_custom.NewLogger().Error(err)
		return
	}
	//init
	id := req.Data.ID
	req.Data.IsMarked = true
	//
	// process update
	unRead := 0
	isUpdate := false
	s.websocketService.Mu.Lock()
	// if topic does not exist, stop the process
	if listNotification, ok := s.websocketService.Notifications[req.Username]; ok {
		for i, item := range listNotification {
			if item.ID == id {
				s.websocketService.Notifications[req.Username][i] = req.Data
				isUpdate = true
				continue
			}
			if !item.IsMarked {
				unRead++
			}
		}
	}
	s.websocketService.Mu.Unlock()

	if !isUpdate {
		err := errors.New("id not found")
		s.websocketService.SendError(conn, err.Error())
		logger_custom.NewLogger().Error(err)
		return
	}

	// update to db
	// TODO
	schema := []entity.UserNotification{
		{
			ID:       id,
			IsMarked: true,
		},
	}
	err := s.notificationRepo.UpdateUserNotification(context.Background(), schema)
	if err != nil {
		return
	}

	// send data updated to user
	notiUpdated := request.MessageRequest{
		Username: username,
		Action:   enum_notification.UpdateAction().Data,
		Data:     []request.Notification{req.Data},
		Unread:   &unRead,
	}

	reqBytes, err := json.Marshal(notiUpdated)
	if err != nil {
		s.websocketService.SendError(conn, err.Error())
		logger_custom.NewLogger().Error(err)
		return
	}
	//sent message
	s.websocketService.Publish(username, reqBytes)
}

// validate
func (s *websocketUsecase) validate(conn *websocket.Conn, req request.MarkAsReadRequest, username string) error {
	// Kiểm tra username và ID
	if req.Username == "" || username == "" || username != req.Username {
		return errors.New("invalid or mismatched username")
	}
	if req.Data.ID == 0 {
		return errors.New("id is empty")
	}

	// Kiểm tra kết nối client
	if clients, ok := s.websocketService.Clients[req.Username]; !ok || clients[conn] == false {
		return errors.New("connection not accepted")
	}

	return nil
}
