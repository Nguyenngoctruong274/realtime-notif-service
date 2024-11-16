package echoprometheus

import (
	"yes4all/ads-noti-api/pkg/prometheus"

	"github.com/labstack/echo/v4"
)

type Middleware interface {
	GetMetrics() []echo.MiddlewareFunc
}
type EmptyMiddleware struct {
}

func (m *EmptyMiddleware) GetMetrics() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type middleware struct {
	serviceName string
	collector   prometheus.Collector
}

func NewMiddleware(serviceName string, collector prometheus.Collector) Middleware {
	m := &middleware{
		serviceName: serviceName,
		collector:   collector,
	}
	m.collectorLatency()
	m.collectorLatencyHistogram()
	m.collectorRequestSize()
	m.collectorRequestTotal()
	return m
}

func (m *middleware) GetMetrics() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{m.Latency, m.LatencyHistogram, m.RequestSize, m.ResponseStatus}
}
