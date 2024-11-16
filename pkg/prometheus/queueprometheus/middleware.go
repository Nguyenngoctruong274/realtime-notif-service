package queueprometheus

import (
	"yes4all/ads-noti-api/pkg/messaging"
	"yes4all/ads-noti-api/pkg/prometheus"
)

type Middleware interface {
	GetMetrics() []messaging.HandlerFunc
}
type EmptyMiddleware struct {
}

func (m *EmptyMiddleware) GetMetrics() []messaging.HandlerFunc {
	return []messaging.HandlerFunc{}
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

func (m *middleware) GetMetrics() []messaging.HandlerFunc {
	return []messaging.HandlerFunc{m.Latency, m.LatencyHistogram, m.RequestSize, m.ResponseStatus}
}
