package router

import (
	"github.com/gin-gonic/gin"
	"yes4all/ads-noti-api/pkg/middleware/auth"
	"yes4all/ads-noti-api/services/ads-noti/http/handler"
)

type ExternalRouter interface {
	Register(routerGroup *gin.RouterGroup)
	RegisterNoMiddleware(r *gin.RouterGroup)
}

type externalRouter struct {
	mwAuth             auth.Auth
	WebsocketHandler   handler.WebsocketHandler
	HealthCheckHandler handler.HealthCheck
}

func NewExternalRouter(
	mwAuth auth.Auth,
	websocketHandler handler.WebsocketHandler,
	healthCheckHandler handler.HealthCheck) ExternalRouter {
	return &externalRouter{
		mwAuth:             mwAuth,
		WebsocketHandler:   websocketHandler,
		HealthCheckHandler: healthCheckHandler}
}

func (sr *externalRouter) Register(r *gin.RouterGroup) {

	g := r.Group("/ws")
	{
		g.POST("/generate-jwt", sr.WebsocketHandler.GenerateJWTHandler)
		g.GET("/connection/:username", sr.mwAuth.Authentication, sr.WebsocketHandler.HandleConnection)
		g.POST("/add-notification", sr.WebsocketHandler.AddNotification)

	}

}

func (sr *externalRouter) RegisterNoMiddleware(r *gin.RouterGroup) {
	healthCheckGroup := r.Group("/health-check")
	{
		healthCheckGroup.GET("/info", sr.HealthCheckHandler.HealthCheckInfo)
	}
}
