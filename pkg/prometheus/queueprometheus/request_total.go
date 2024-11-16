package queueprometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"yes4all/ads-noti-api/pkg/messaging"
	"yes4all/ads-noti-api/pkg/xerrors"
)

var (
	totalRequestCounterOpts = prometheus.CounterOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "queue_metric_requests_total",
		Help:      "Total number requests by queue name and error code",
	}
	totalRequestCollector = prometheus.NewCounterVec(totalRequestCounterOpts,
		[]string{"event", "key", "error_code"})
)

func (m *middleware) collectorRequestTotal() {
	totalRequestCounterOpts.Subsystem = m.serviceName
	totalRequestCollector = prometheus.NewCounterVec(totalRequestCounterOpts,
		[]string{"event", "key", "error_code"})
	m.collector.Collect(totalRequestCollector)
}

func (m *middleware) ResponseStatus(c *messaging.Context) {
	c.Next()
	go m.sendResponseStatus(c)
}

func (m *middleware) sendResponseStatus(c *messaging.Context) {
	statusCode := xerrors.Success
	if c.Err() != nil {
		statusCode = xerrors.GetErrorType(c.Err())
	}

	totalRequestCollector.WithLabelValues(c.QueueName(), c.Key(), fmt.Sprint(statusCode)).Inc()
}
