package echoprometheus

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	reqSizeGaugeOpts = prometheus.GaugeOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_handler_request_size",
		Help:      "Request size.",
	}
	requestSizeCollector = prometheus.NewGaugeVec(reqSizeGaugeOpts, []string{"url", "method"})
)

func (m *middleware) collectorRequestSize() {
	latencyGaugeOpts.Subsystem = m.serviceName
	requestSizeCollector = prometheus.NewGaugeVec(reqSizeGaugeOpts, []string{"url", "method"})
	m.collector.Collect(requestSizeCollector)
}

func (m *middleware) RequestSize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		go m.sendRequestSize(c)
		return next(c)
	}
}

func (m *middleware) sendRequestSize(c echo.Context) {
	reqSize := computeApproximateRequestSize(c.Request())
	requestSizeCollector.WithLabelValues(c.Path(), c.Request().Method).Set(float64(reqSize))
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
