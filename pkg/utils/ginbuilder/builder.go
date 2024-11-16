package ginbuilder

import (
	"yes4all/ads-noti-api/pkg/utils/ginutils"

	"github.com/gin-gonic/gin"
)

type builder struct {
	middlewares []gin.HandlerFunc
}

func BaseBuilder() *builder {
	return &builder{
		middlewares: []gin.HandlerFunc{
			gin.Recovery(),
		},
	}
}

func Default() *builder {
	return &builder{}
}

func (b *builder) WithBodyLogger(skipPaths ...string) *builder {
	b.middlewares = append(b.middlewares, ginutils.Logger(skipPaths...))
	return b
}

func (b *builder) WithTraceID() *builder {
	b.middlewares = append(b.middlewares, ginutils.InjectTraceID)
	return b
}

func (b *builder) Build() *gin.Engine {
	e := defaultGinEngine()
	e.Use(b.middlewares...)
	return e
}

func defaultGinEngine() *gin.Engine {
	e := gin.New()
	return e
}
