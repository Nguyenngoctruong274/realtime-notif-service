package echoprometheus

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	latencyGaugeOpts = prometheus.GaugeOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_handler_latency",
		Help:      "Request latencies by URI & HTTP status code.",
	}
	latencyCollector = prometheus.NewGaugeVec(latencyGaugeOpts, []string{"url", "method", "http_status_code"})
)

func (m *middleware) collectorLatency() {
	latencyGaugeOpts.Subsystem = m.serviceName
	latencyCollector = prometheus.NewGaugeVec(latencyGaugeOpts, []string{"url", "method", "http_status_code"})
	m.collector.Collect(latencyCollector)
}

func (m *middleware) Latency(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		go m.sendLatency(c, start)
		return nil
	}
}

func (m *middleware) sendLatency(c echo.Context, start time.Time) {
	statusCode := c.Response().Status
	httpStatusCode := fmt.Sprintf("%dxx", statusCode/statusCodeDivisor)
	latency := time.Since(start).Seconds()

	latencyCollector.WithLabelValues(c.Path(), c.Request().Method, httpStatusCode).Set(latency)
}
