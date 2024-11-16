package queueprometheus

import (
	"fmt"
	"time"
	"yes4all/ads-noti-api/pkg/messaging"
	"yes4all/ads-noti-api/pkg/xerrors"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	latencyGaugeOpts = prometheus.GaugeOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "queue_metric_latencies",
		Help:      "Message latencies",
	}
	latencyCollector = prometheus.NewGaugeVec(latencyGaugeOpts, []string{"event", "key", "error_code"})
)

func (m *middleware) collectorLatency() {
	latencyGaugeOpts.Subsystem = m.serviceName
	latencyCollector = prometheus.NewGaugeVec(latencyGaugeOpts, []string{"event", "key", "error_code"})
	m.collector.Collect(latencyCollector)
}

func (m *middleware) Latency(c *messaging.Context) {
	start := time.Now()
	c.Next()
	go m.sendLatency(c, start)
}

func (m *middleware) sendLatency(c *messaging.Context, start time.Time) {
	latency := time.Since(start).Seconds()
	statusCode := xerrors.Success
	if c.Err() != nil {
		statusCode = xerrors.GetErrorType(c.Err())
	}
	latencyCollector.WithLabelValues(c.QueueName(), c.Key(), fmt.Sprint(statusCode)).Set(latency)
}
