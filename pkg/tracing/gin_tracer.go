package tracing

import (
	"fmt"
	"yes4all/ads-noti-api/pkg/utils/ginutils"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Tracer interface {
	TracingHandler(c *gin.Context)
}
type tracer struct {
	tracer opentracing.Tracer
}

func NewTracer(t opentracing.Tracer) *tracer {
	return &tracer{t}
}

func (t *tracer) TracingHandler(c *gin.Context) {
	spanCtx, _ := t.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	opsName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
	serverSpan, ctx := StartSpanFromCtx(c, opsName, ext.RPCServerOption(spanCtx))
	ext.HTTPUrl.Set(serverSpan, c.FullPath())
	ext.HTTPMethod.Set(serverSpan, c.Request.Method)
	c.Request = c.Request.WithContext(ctx)
	defer func() {
		ext.HTTPStatusCode.Set(serverSpan, uint16(c.Writer.Status()))
		serverSpan.SetTag("request_id", ginutils.GetTraceIDFromCtx(c))
		if c.Errors != nil {
			serverSpan.SetTag("errors", c.Errors)
		}
		serverSpan.Finish()
	}()
	c.Next()
}
