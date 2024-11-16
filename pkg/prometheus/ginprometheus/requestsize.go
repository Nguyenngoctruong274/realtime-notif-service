package ginprometheus

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	reqSizeGaugeOpts = prometheus.GaugeOpts{
		Namespace: defaultNamespace,
		Subsystem: defaultSubSystem,
		Name:      "http_metric_request_size",
		Help:      "Request size.",
	}
	requestSizeCollector = prometheus.NewGaugeVec(reqSizeGaugeOpts, []string{"url", "method"})
)

func (m *middleware) collectorRequestSize() {
	reqSizeGaugeOpts.Subsystem = m.serviceName
	requestSizeCollector = prometheus.NewGaugeVec(reqSizeGaugeOpts, []string{"url", "method"})
	m.collector.Collect(requestSizeCollector)
}

func (m *middleware) RequestSize(c *gin.Context) {
	go m.sendRequestSize(c)
	c.Next()
}

func (m *middleware) sendRequestSize(c *gin.Context) {
	reqSize := computeApproximateRequestSize(c.Request)
	requestSizeCollector.WithLabelValues(c.FullPath(), c.Request.Method).Set(float64(reqSize))
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
