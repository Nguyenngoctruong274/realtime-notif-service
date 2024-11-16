package ginprometheus

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	latencyGaugeOpts = prometheus.GaugeOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_latencies",
		Help:      "Request latency by URI & HTTP status code.",
	}
	latencyCollector = prometheus.NewGaugeVec(latencyGaugeOpts, []string{"url", "method", "http_status_code"})
)

func (m *middleware) collectorLatency() {
	latencyGaugeOpts.Subsystem = m.serviceName
	latencyCollector = prometheus.NewGaugeVec(latencyGaugeOpts, []string{"url", "method", "http_status_code"})
	m.collector.Collect(latencyCollector)
}

func (m *middleware) Latency(c *gin.Context) {
	start := time.Now()
	c.Next()
	go m.sendLatency(c, start)
}

func (m *middleware) sendLatency(c *gin.Context, start time.Time) {
	statusCode := c.Writer.Status()
	httpStatusCode := fmt.Sprintf("%dxx", statusCode/statusCodeDivisor)
	latency := time.Since(start).Seconds()

	latencyCollector.WithLabelValues(c.FullPath(), c.Request.Method, httpStatusCode).Set(latency)
}
