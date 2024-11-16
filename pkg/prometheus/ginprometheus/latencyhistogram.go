package ginprometheus

import (
	"time"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	latencyHistogramOpts = prometheus.HistogramOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_request_duration",
		Help:      "Request duration histogram in seconds",
		Buckets:   prometheus.ExponentialBuckets(bucketStart, bucketFactor, bucketCount),
	}
	latencyHistogramCollector = prometheus.NewHistogramVec(
		latencyHistogramOpts, []string{"url", "method"})
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
	latencyHistogramCollector = prometheus.NewHistogramVec(latencyHistogramOpts, []string{"url", "method"})
	m.collector.Collect(latencyHistogramCollector)
}

func (m *middleware) LatencyHistogram(c *gin.Context) {
	start := time.Now()
	c.Next()
	go m.sendLatencyHistogram(c, start)
}
func (m *middleware) sendLatencyHistogram(c *gin.Context, start time.Time) {
	latencyHistogramCollector.WithLabelValues(c.FullPath(), c.Request.Method).Observe(time.Since(start).Seconds())
}
