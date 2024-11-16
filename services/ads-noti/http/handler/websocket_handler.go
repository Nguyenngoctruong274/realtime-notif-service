package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"yes4all/ads-noti-api/pkg/middleware/auth"
	"yes4all/ads-noti-api/pkg/utils/common"
	"yes4all/ads-noti-api/pkg/xcontext"
	"yes4all/ads-noti-api/services/ads-noti/model/request"
	"yes4all/ads-noti-api/services/ads-noti/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Kiểm tra Origin
	},
}

type WebsocketHandler interface {
	GenerateJWTHandler(c *gin.Context)
	HandleConnection(c *gin.Context)
	AddNotification(c *gin.Context)
}

type websocketHandler struct {
	Auth             auth.Auth
	WebsocketUsecase usecase.WebsocketUsecase
}

func NewNotificationManagementHandler(
	auth auth.Auth,
	websocketUsecase usecase.WebsocketUsecase) WebsocketHandler {
	return &websocketHandler{
		Auth:             auth,
		WebsocketUsecase: websocketUsecase,
	}
}

// GET Handle Connection godoc
// @Summary GET Handle Connection
// @Description GET Handel Connection
// @Tags cyborg
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param AProfileId header string true "AProfileId"
// @Param request body request.GetProfilesByAsinReq true "model"
// @Success 200 {object} common.Response{data=response.GetProfilesByAsinResp}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /external/ws/connection/{username} [Get]
func (n *websocketHandler) HandleConnection(c *gin.Context) {
	ctx := xcontext.CloneByGinContext(c)
	username := c.Param("username")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		err = fmt.Errorf("error upgrading to websocket_service: %s", err)
		return
	}
	err = n.WebsocketUsecase.HandleConnection(ctx, conn, username)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	common.JSONOk(c, nil)
}

// // POST Add Notification Websocket godoc
// // @Summary POST Add Notification Websocket
// // @Description POST Add Notification Websocket
// // @Tags cyborg
// // @Accept json
// // @Produce json
// // @Security ApiKeyAuth
// // @param Authorization header string true "Authorization"
// // @param AProfileId header string true "AProfileId"
// // @Param request body request.GetProfilesByAsinReq true "model"
// // @Success 200 {object} common.Response
// // @Failure 400 {object} common.Response
// // @Failure 500 {object} common.Response
// // @Router /external/ws/add-notification [POST]
func (n *websocketHandler) AddNotification(c *gin.Context) {
	req := request.MessageRequest{}
	if err := c.BindJSON(&req); err != nil {
		common.HandleBindError(c, err, "AddNotification")
		return
	}
	ctx := xcontext.CloneByGinContext(c)
	err := n.WebsocketUsecase.AddNotification(ctx, req)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	common.JSONOk(c, nil)
}

// Post Generate JWT godoc
// @Summary Generate JWT Connection
// @Description Generate JWT Connection
// @Tags cyborg
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param AProfileId header string true "AProfileId"
// @Param request body request.JWTRequest true "model"
// @Success 200 {object} common.Response{data=response.JWTResp}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /external/ws/generate-jwt [Post]
func (n *websocketHandler) GenerateJWTHandler(c *gin.Context) {
	ctx := xcontext.CloneByGinContext(c)
	req := request.JWTRequest{}
	if err := c.BindJSON(&req); err != nil {
		common.HandleBindError(c, err, "GenerateJWTHandler")
		return
	}
	data, err := n.Auth.GenerateJWTHandler(ctx, req)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	common.JSONOk(c, data)
}
