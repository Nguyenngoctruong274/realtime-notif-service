package tracingfx

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"

	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/tracing"
)

var Module = fx.Provide(provideTracer, provideTracerHandler)

func provideTracer(lifecycle fx.Lifecycle) (opentracing.Tracer, error) {
	cfg := config.TracingConfig()
	tracer, closer, err := tracing.Initialize(cfg)
	if err != nil {
		return nil, err
	}
	lifecycle.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		return closer.Close()
	}})
	return tracer, nil
}

func provideTracerHandler(tracer opentracing.Tracer) tracing.Tracer {
	return tracing.NewTracer(tracer)
}
