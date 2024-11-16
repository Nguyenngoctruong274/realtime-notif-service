package tracing

import (
	"context"
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	appconfig "yes4all/ads-noti-api/pkg/config"
)

func Initialize(cfg appconfig.JaegerConfig) (tracer opentracing.Tracer, closer io.Closer, err error) {
	jaegerConfig, err := config.FromEnv()
	if err != nil {
		log.Fatal("Failed to read .env file for config jaeger agent.")
	}

	tracer, closer, err = jaegerConfig.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}

func StartSpanFromCtx(ctx context.Context,
	operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}
