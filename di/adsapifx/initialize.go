package adsapifx

import (
	"yes4all/ads-noti-api/pkg/middleware/auth"
	"yes4all/ads-noti-api/services/ads-noti/http/handler"
	"yes4all/ads-noti-api/services/ads-noti/http/router"
	"yes4all/ads-noti-api/services/ads-noti/repository"
	"yes4all/ads-noti-api/services/ads-noti/usecase"

	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Provide(
	// middleware auth cloak
	// Routers
	provideExternalRouter,
	// HandlersInternal
	provideNotificationManagementHandler,

	// Repository
	provideNotificationRepo,

	// HandlersExternal
	provideHealthCheckHandler,
	provideAuth,
	// Cache

	// Usercases
	provideWebsocketUsecase,
	//usecase common

	// HTTP Client

	// authCloak

)

// ======== Router ===========//
func provideExternalRouter(mwAuth auth.Auth, websocketHandler handler.WebsocketHandler, healthCheckHandler handler.HealthCheck) router.ExternalRouter {
	return router.NewExternalRouter(
		mwAuth,
		websocketHandler, healthCheckHandler)
}

// ======== Handler ===========//
func provideNotificationManagementHandler(
	auth auth.Auth,
	websocketUsecase usecase.WebsocketUsecase) handler.WebsocketHandler {
	return handler.NewNotificationManagementHandler(auth, websocketUsecase)
}
func provideHealthCheckHandler() handler.HealthCheck {
	return handler.NewHealthCheck()
}

func provideAuth() auth.Auth {
	return auth.NewAuth()
}

// ======== Repository ===========//
func provideNotificationRepo() repository.NotificationRepo {
	return repository.NewNotificationRepo()
}

// ======== Usecase ===========//
func provideWebsocketUsecase() usecase.WebsocketUsecase {
	return usecase.NewWebsocketUsecase(repository.NewNotificationRepo())
}
