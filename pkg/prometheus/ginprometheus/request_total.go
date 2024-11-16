package ginprometheus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	totalRequestCounterOpts = prometheus.CounterOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_requests_total",
		Help:      "Total number requests by URI & HTTP status code.",
	}
	totalRequestCollector = prometheus.NewCounterVec(totalRequestCounterOpts,
		[]string{"url", "method", "http_status_code"})
)

func (m *middleware) collectorRequestTotal() {
	totalRequestCounterOpts.Subsystem = m.serviceName
	totalRequestCollector = prometheus.NewCounterVec(totalRequestCounterOpts,
		[]string{"url", "method", "http_status_code"})
	m.collector.Collect(totalRequestCollector)
}

func (m *middleware) ResponseStatus(c *gin.Context) {
	c.Next()
	go m.sendResponseStatus(c)
}

func (m *middleware) sendResponseStatus(c *gin.Context) {
	statusCode := c.Writer.Status()
	httpStatusCode := fmt.Sprintf("%dxx", statusCode/statusCodeDivisor)

	totalRequestCollector.WithLabelValues(c.FullPath(), c.Request.Method,
		httpStatusCode).Inc()
}
