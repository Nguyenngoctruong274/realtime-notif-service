package tracing

import (
	"context"
	"fmt"

	"go.elastic.co/apm/transport"

	"go.elastic.co/apm"
)

func Initialize() (tracer *apm.Tracer, err error) {
	apmTransport, err := transport.NewHTTPTransport()
	if err != nil {
		return nil, err
	}
	return apm.NewTracerOptions(apm.TracerOptions{Transport: apmTransport})
}

func CloseDefaultTracer() {
	if apm.DefaultTracer != nil {
		apm.DefaultTracer.Close()
	}
}

func StartSpan(ctx context.Context, spanName string) (*apm.Span, context.Context) {
	return StartSpanWithType(ctx, spanName, SpanTypeUnknown)
}

func StartSpanWithType(ctx context.Context, spanName string, spanType SpanType) (*apm.Span, context.Context) {
	return apm.StartSpan(ctx, spanName, fmt.Sprint(spanType))
}
