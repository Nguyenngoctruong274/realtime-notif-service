package echoprometheus

import (
	"time"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/labstack/echo/v4"

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

func (m *middleware) LatencyHistogram(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		go m.sendLatencyHistogram(c, start)
		return nil
	}
}

func (m *middleware) sendLatencyHistogram(c echo.Context, start time.Time) {
	latencyHistogramCollector.WithLabelValues(c.Path(), c.Request().Method).Observe(time.Since(start).Seconds())
}
