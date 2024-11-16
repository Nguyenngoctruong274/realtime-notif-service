package swagger

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"yes4all/ads-noti-api/docs"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/utils/constants"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type swagger struct {
}

func NewSwagger() *swagger {
	return &swagger{}
}

func (m *swagger) Register(gGroup gin.IRouter) {
	g := gGroup.Group("")
	{
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
		if config.ServerConfig().Env == constants.DevelopmentEnv {
			docs.SwaggerInfo.Schemes = []string{"http"}
		}
		g.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (m *swagger) SwaggerHandler(isProduction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProduction {
			return
		}
		docs.SwaggerInfo.Title = "ads api service swagger"
		docs.SwaggerInfo.Description = "Thông tin các api của ads api service"
		docs.SwaggerInfo.Host = strings.ToLower(c.Request.Host)
		if os.Getenv("ENVIRONMENT") != "local" {
			docs.SwaggerInfo.Host = "yams.test.yes4all.internal/ads-noti"
		}
		docs.SwaggerInfo.BasePath = "/api"
		c.Next()
	}
}
