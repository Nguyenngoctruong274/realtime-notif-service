package cmd

import (
	"context"
	"fmt"
	"os"
	"time"
	"yes4all/ads-noti-api/cmd/banner"
	"yes4all/ads-noti-api/di/adsapifx"
	"yes4all/ads-noti-api/di/metricsfx"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/difx/prometheusfx"
	"yes4all/ads-noti-api/pkg/difx/tracingfx"
	"yes4all/ads-noti-api/pkg/graceful"
	"yes4all/ads-noti-api/pkg/infra"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/queue/kafka"
	"yes4all/ads-noti-api/pkg/swagger"
	"yes4all/ads-noti-api/pkg/tracing"
	"yes4all/ads-noti-api/pkg/utils/common"
	"yes4all/ads-noti-api/pkg/utils/constants"
	"yes4all/ads-noti-api/pkg/utils/errors"
	"yes4all/ads-noti-api/pkg/utils/ginbuilder"
	"yes4all/ads-noti-api/pkg/utils/ginutils"
	"yes4all/ads-noti-api/pkg/websocket_service"
	"yes4all/ads-noti-api/services/ads-noti/http/router"
	"yes4all/ads-noti-api/services/ads-noti/kafka/notification"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Somethings short",
	Long:  "Somethings long",
	Run: func(_ *cobra.Command, _ []string) {
		NewServer().Run()
	},
	Version: "1.0.0",
}

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	app := fx.New(
		fx.Invoke(logger.InitLogger),
		fx.Invoke(errors.Initialize),
		fx.Invoke(config.InitConfig),
		fx.Invoke(infra.InitPostgresql),
		fx.Invoke(websocket_service.InitWebsocket),
		// fx.Invoke(infra.InitPostgresqlReport),

		tracingfx.Module,
		metricsfx.Module,

		prometheusfx.Module,
		adsapifx.Module,
		kafka.Module,
		notification.InitKafkaConn,
		notification.SubsribeToTopic,
		notification.SubsribeToTopicDLQ,
		fx.Provide(provideGinEngine),
		fx.Invoke(
			registerSwaggerHandler,
			registerService),
		fx.Invoke(startServer),
		fx.Invoke(banner.Print),
	)
	app.Run()
}

func provideGinEngine() *gin.Engine {
	return ginbuilder.BaseBuilder().
		WithBodyLogger("/metrics", "/swagger/*any").
		Build()
}

func registerSwaggerHandler(g *gin.Engine) {
	swaggerAPI := g.Group("/swagger")
	swag := swagger.NewSwagger()
	swaggerAPI.Use(swag.SwaggerHandler(config.ServerConfig().Env == constants.ProductionEnv))
	swag.Register(swaggerAPI)
}

func registerService(g *gin.Engine,
	tracer tracing.Tracer,
	externalRouter router.ExternalRouter,
) {
	// 	add_header 'Access-Control-Allow-Origin' * always;
	// add_header 'Access-Control-Allow-Credentials' 'true' always;
	// add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, DELETE, PUT' always;
	// add_header 'Access-Control-Allow-Headers' * always;
	// add_header 'Access-Control-Max-Age' 1728000;
	// if ($request_method = 'OPTIONS') {
	//     return 204;
	// }
	// ~
	// set cors prolicy for local testing frontend
	if os.Getenv("ENVIRONMENT") == "local" {
		g.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "OPTIONS", "GET"},
			AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
			ExposeHeaders:    []string{"Content-Length", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	g.NoRoute(func(c *gin.Context) {
		common.HandlePathNotExistError(c, errors.ErrNotExistPath)
	})
	g.Use(ginutils.RecoverPanic)
	// Thiết lập route cho WebSocket

	// group API internal
	externalGroup := g.Group("/api/external")
	externalGroup.Use(ginutils.InjectTraceID, tracer.TracingHandler)
	//externalGroup.Use(auth.VerifyAuth)
	externalRouter.Register(externalGroup)
	externalRouter.RegisterNoMiddleware(externalGroup)

}

func startServer(lifecycle fx.Lifecycle, g *gin.Engine) {
	gracefulService := graceful.NewService(graceful.WithStopTimeout(time.Second), graceful.WithWaitTime(time.Second))
	gracefulService.Register(g)
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := fmt.Sprintf("%d", config.ServerConfig().HTTPPort)
				fmt.Println("run on port:", port)
				go gracefulService.StartServer(g, port)
				return nil
			},
			OnStop: func(context.Context) error {
				gracefulService.Close()
				infra.ClosePostgresql()
				websocket_service.CloseServer() //nolint

				//infra.CloseRedis()

				return nil
			},
		},
	)
}
