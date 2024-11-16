package queueprometheus

import (
	"yes4all/ads-noti-api/pkg/messaging"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	reqSizeGaugeOpts = prometheus.GaugeOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "queue_metric_request_size",
		Help:      "Request size.",
	}
	requestSizeCollector = prometheus.NewGaugeVec(reqSizeGaugeOpts, []string{"event", "key"})
)

func (m *middleware) collectorRequestSize() {
	reqSizeGaugeOpts.Subsystem = m.serviceName
	requestSizeCollector = prometheus.NewGaugeVec(reqSizeGaugeOpts, []string{"event", "key"})
	m.collector.Collect(requestSizeCollector)
}

func (m *middleware) RequestSize(c *messaging.Context) {
	go m.sendRequestSize(c)
	c.Next()
}

func (m *middleware) sendRequestSize(c *messaging.Context) {
	reqSize := computeApproximateRequestSize(c)
	requestSizeCollector.WithLabelValues(c.QueueName(), c.Key()).Set(float64(reqSize))
}

func computeApproximateRequestSize(c *messaging.Context) int {
	s := 0
	for name, values := range c.Headers() {
		s += len(name)
		s += len(values)
	}
	s += len(c.Data())
	return s
}
