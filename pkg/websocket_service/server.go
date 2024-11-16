package websocket_service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	logger_custom "yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/services/ads-noti/model/enum/enum_notification"
	"yes4all/ads-noti-api/services/ads-noti/model/request"
)

var server *Server

type Server struct {
	Clients       Clients
	Notifications map[string][]request.Notification
	Mu            sync.Mutex
	IsConnect     map[string]bool
}

func InitWebsocket() {
	var once sync.Once
	once.Do(func() {
		server = &Server{
			IsConnect:     make(map[string]bool),
			Clients:       make(Clients),
			Notifications: make(map[string][]request.Notification),
			Mu:            sync.Mutex{},
		}
	})
}

func GetServerSocket() *Server {
	if server == nil {
		logger_custom.NewLogger().Fatal("Connection websocket_service failed")
	}
	return server
}

func CloseServer() {
	// Đóng tất cả các kết nối Client
	for username, client := range server.Clients {
		for conn, ok := range client {
			if ok {
				conn.Close() // Đóng kết nối WebSocket
			}
		}
		delete(server.Clients, username)
		delete(server.IsConnect, username)
	}
	// Xóa thông tin lưu trữ các thông báo toàn cục
	server.Notifications = nil
}

// Publish sends a message to all subscribing clients of a topic
func (s *Server) Publish(username string, message []byte) {

	// if client does not exist, stop the process
	if _, exist := s.Clients[username]; !exist {
		return
	}
	// if client exist
	client := s.Clients[username]
	// send the message to the clients
	var wg sync.WaitGroup
	for conn := range client {
		// add 1 job to wait group
		wg.Add(1)
		// send with goroutines
		go s.SendWithWait(conn, message, &wg)
	}
	// wait until all goroutines jobs done
	wg.Wait()
}

// Subscribe adds a client to a topic's client map
func (s *Server) Subscribe(conn *websocket.Conn, username string) *Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	//save first connect user
	if !s.IsConnect[username] {
		s.IsConnect[username] = true
	}

	// if topic exist, check the client map
	if connections, exist := s.Clients[username]; exist {
		// if client already subbed, stop the process
		if _, subbed := connections[conn]; !subbed {
			connections[conn] = true
		}
		// if not subbed, add to client map
	} else {
		// if topic does not exist, create a new topic
		s.Clients[username] = make(map[*websocket.Conn]bool)
		// add the client to the topic
		s.Clients[username][conn] = true
	}

	//save data
	if notifications, exists := s.Notifications[username]; exists {
		count := 0
		for _, notification := range notifications {
			if !notification.IsMarked {
				count++
			}
		}
		data := request.MessageRequest{
			Username: username,
			Action:   enum_notification.ListAction().Data,
			Data:     notifications,
			Unread:   &count,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			s.SendError(conn, errInvalidMessage)
			return s
		}
		s.Send(conn, string(dataBytes))
	}
	return s
}

// Unsubscribe removes a clients from a topic's client map
func (s *Server) Unsubscribe(conn *websocket.Conn, username string) *Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if connections, ok := s.Clients[username]; ok {
		if _, exists := connections[conn]; exists {
			delete(connections, conn)
		}
		// Nếu client không còn kết nối nào, xóa luôn client khỏi map
		if len(connections) == 0 {
			delete(s.Clients, username)
		}
	}
	return s
}

// Send simply sends message to the websocket_service client
func (s *Server) Send(conn *websocket.Conn, message string) {
	// send simple message
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}
func (s *Server) SendError(conn *websocket.Conn, message string) {
	// send simple message
	messErr := request.MessageRequest{
		Username: "",
		Action:   enum_notification.ErrorAction().Data,
		Data:     nil,
		Message:  &message,
		Unread:   nil,
	}
	notiRequest, _ := json.Marshal(messErr)
	conn.WriteMessage(websocket.TextMessage, notiRequest)
}

// SendWithWait sends message to the websocket_service client using wait group, allowing usage with goroutines
func (s *Server) SendWithWait(conn *websocket.Conn, message []byte, wg *sync.WaitGroup) {
	// send simple message
	conn.WriteMessage(websocket.TextMessage, message)

	// set the task as done
	wg.Done()
}
