package ginprometheus

import (
	"fmt"
	"net/http"
	"time"
)

type MetricData struct {
	Latency      func(r *http.Request, start time.Time, statusCode int, groupPath string)
	Histogram    func(r *http.Request, start time.Time, groupPath string)
	RequestTotal func(r *http.Request, statusCode int, groupPath string)
	RequestSize  func(r *http.Request, groupPath string)
	ResetLatency func()
}

var prometheusMetric *MetricData

func InitPrometheusMetric() {
	prometheusMetric = &MetricData{
		Latency:      sendLatencyMetric,
		Histogram:    sendLatencyHistogramMetric,
		RequestTotal: sendResponseStatusMetric,
		RequestSize:  sendRequestSizeMetric,
		ResetLatency: resetLatencyMetric,
	}
}

func GetPrometheusMetric() *MetricData {
	return prometheusMetric
}

func sendLatencyMetric(r *http.Request, start time.Time, statusCode int, groupPath string) {
	httpStatusCode := fmt.Sprintf("%dxx", statusCode/statusCodeDivisor)
	latency := time.Since(start).Seconds()

	tempURL := groupPath
	if tempURL == "" {
		tempURL = r.URL.Path
	}

	latencyCollector.WithLabelValues(tempURL, r.Method, httpStatusCode).Set(latency)
}

func sendLatencyHistogramMetric(r *http.Request, start time.Time, groupPath string) {
	tempURL := groupPath
	if tempURL == "" {
		tempURL = r.URL.Path
	}

	latencyHistogramCollector.WithLabelValues(tempURL, r.Method).Observe(time.Since(start).Seconds())
}

func sendResponseStatusMetric(r *http.Request, statusCode int, groupPath string) {
	tempURL := groupPath
	if tempURL == "" {
		tempURL = r.URL.Path
	}

	httpStatusCode := fmt.Sprintf("%dxx", statusCode/statusCodeDivisor)
	totalRequestCollector.WithLabelValues(tempURL, r.Method, httpStatusCode).Inc()
}

func sendRequestSizeMetric(r *http.Request, groupPath string) {
	tempURL := groupPath
	if tempURL == "" {
		tempURL = r.URL.Path
	}

	reqSize := computeRequestSize(r)
	requestSizeCollector.WithLabelValues(tempURL, r.Method).Set(float64(reqSize))
}

func computeRequestSize(r *http.Request) int {
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

func resetLatencyMetric() {
	latencyCollector.Reset()
}
