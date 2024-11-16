package queueprometheus

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/messaging"
	"yes4all/ads-noti-api/pkg/xerrors"
)

var (
	latencyHistogramOpts = prometheus.HistogramOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "queue_metric_request_duration",
		Help:      "Request duration histogram in seconds",
		Buckets:   prometheus.ExponentialBuckets(bucketStart, bucketFactor, bucketCount),
	}
	latencyHistogramCollector = prometheus.NewHistogramVec(
		latencyHistogramOpts, []string{"event", "key", "error_code"})
)

const (
	bucketStart  = 0.01
	bucketFactor = 2
	bucketCount  = 10
)

func (m *middleware) collectorLatencyHistogram() {
	logger.NewLogger().Info("LinearBuckets creates 'count' buckets, each 'width' wide,"+
		" where the lowest bucket has an upper bound of 'start'",
		"count", bucketCount, "Factor", bucketFactor, "start", bucketStart)
	latencyHistogramOpts.Subsystem = m.serviceName
	latencyHistogramCollector = prometheus.NewHistogramVec(latencyHistogramOpts, []string{"event", "key", "error_code"})
	m.collector.Collect(latencyHistogramCollector)
}

func (m *middleware) LatencyHistogram(c *messaging.Context) {
	start := time.Now()
	c.Next()
	go m.sendLatencyHistogram(c, start)
}

func (m *middleware) sendLatencyHistogram(c *messaging.Context, start time.Time) {
	statusCode := xerrors.Success
	if c.Err() != nil {
		statusCode = xerrors.GetErrorType(c.Err())
	}
	latency := time.Since(start).Seconds()
	latencyHistogramCollector.WithLabelValues(c.QueueName(), c.Key(), fmt.Sprint(statusCode)).Observe(latency)
}
