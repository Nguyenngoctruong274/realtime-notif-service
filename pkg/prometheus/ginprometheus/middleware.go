package ginprometheus

import (
	"github.com/gin-gonic/gin"
	"yes4all/ads-noti-api/pkg/prometheus"
)

type Middleware interface {
	GetMetrics() []gin.HandlerFunc
}
type EmptyMiddleware struct {
}

func (m *EmptyMiddleware) SetMetrics(c *gin.Context) {
}

func (m *EmptyMiddleware) GetMetrics() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
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

func (m *middleware) GetMetrics() []gin.HandlerFunc {
	return []gin.HandlerFunc{m.Latency, m.LatencyHistogram, m.RequestSize, m.ResponseStatus}
}
