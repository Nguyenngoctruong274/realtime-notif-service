package prompusher

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/xerrors"
)

var (
	promLatencyCollector = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prom_http_push_metric_latency",
			Help: "Request latencies by URI & HTTP status code.",
		}, []string{"url", "method", "code"})

	promRequestBodySizeCollector = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prom_http_push_metric_body_size",
			Help: "Request body size by URI & HTTP status code.",
		}, []string{"url", "method", "code"})

	promRequestCounterCollector = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "prom_http_pusher_metric_total",
			Help: "Push metric total",
		}, []string{"url", "method", "code"})
)

type HTTPClient struct {
	client *http.Client
	config GatewayConfig

	lastRequest   int64
	hasSkippedReq int32
}

func NewHTTPClient(cfg GatewayConfig) *HTTPClient {
	client := &http.Client{Transport: http.DefaultTransport, Timeout: cfg.HTTPTimeout}
	return &HTTPClient{
		client: client,
		config: cfg,
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	now := time.Now()
	previous := time.Unix(c.lastRequest, 0)
	if previous.Add(c.config.RequestInterval).Before(now) ||
		!atomic.CompareAndSwapInt64(&c.lastRequest, c.lastRequest, now.Unix()) {
		c.checkRetry(req, now)

		return nil, xerrors.Success.New()
	}

	return c.sendRequest(req, now)
}

func (c *HTTPClient) checkRetry(req *http.Request, start time.Time) {
	if atomic.CompareAndSwapInt32(&c.hasSkippedReq, 0, 1) {
		go c.retrySendRequest(req, start)
	}
}

func (c *HTTPClient) retrySendRequest(req *http.Request, start time.Time) {
	time.Sleep(c.config.RequestInterval)
	c.hasSkippedReq = 0
	nowUnix := start.Unix()
	if nowUnix > c.lastRequest {
		resp, err := c.sendRequest(req, start)
		if err != nil {
			logger.NewLogger().WithKeyword(req.Context(), "failed to retry send prom").WithError(err).Error()
			return
		}
		defer resp.Body.Close() //nolint
	}
}

func (c *HTTPClient) sendRequest(req *http.Request, start time.Time) (*http.Response, error) {
	resp, err := c.client.Do(req)

	go c.updateMetric(req, resp, err, start)

	return resp, err
}

func (c *HTTPClient) updateMetric(req *http.Request,
	resp *http.Response, err error, startTime time.Time) {
	url := req.URL.Path
	method := req.Method
	contentLength := req.ContentLength

	var code string
	if err != nil {
		code = "-1"
	} else if resp != nil {
		code = fmt.Sprintf("%d", resp.StatusCode)
	}

	latency := time.Since(startTime).Seconds()

	promLatencyCollector.WithLabelValues(url, method, code).
		Set(latency)

	promRequestBodySizeCollector.WithLabelValues(url, method, code).
		Set(float64(contentLength))

	promRequestCounterCollector.WithLabelValues(url, method, code).
		Inc()
}
