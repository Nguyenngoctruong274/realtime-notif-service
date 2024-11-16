package echoprometheus

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	totalRequestCounterOpts = prometheus.CounterOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_handler_requests_total",
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

func (m *middleware) ResponseStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		go m.sendResponseStatus(c)
		return nil
	}
}

func (m *middleware) sendResponseStatus(c echo.Context) {
	statusCode := c.Response().Status
	httpStatusCode := fmt.Sprintf("%dxx", statusCode/statusCodeDivisor)

	totalRequestCollector.WithLabelValues(c.Path(), c.Request().Method,
		httpStatusCode).Inc()
}
