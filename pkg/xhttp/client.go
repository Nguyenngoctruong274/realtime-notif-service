package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/prometheus/ginprometheus"
	"yes4all/ads-noti-api/pkg/utils/constants"
	"yes4all/ads-noti-api/pkg/utils/timeutils"
	"yes4all/ads-noti-api/pkg/utils/tracking"
	"yes4all/ads-noti-api/pkg/xcontext"

	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
)

const (
	defaultTimeout       = 30 * time.Second
	defaultLogBodyLength = 3000
	defaultNamespace     = "yams"
	defaultSubsystem     = "yams"
)

// nolint: lll
// Không cần check long line linter cho interface
type Client interface {
	PostJSONWithCustomHeader(ctx context.Context, url string, data, target interface{}, customHeader http.Header, method string, reqOpts ...RequestOption) (statusCode int, err error)
	PostJSON(c context.Context, url string, data, target interface{}, reqOptions ...RequestOption) (int, error)
	PostForm(c context.Context, url string, data, target interface{}, reqOptions ...RequestOption) (int, error)
	Get(c context.Context, url string, target interface{}, reqOptions ...RequestOption) (int, error)
	GetWithQuery(c context.Context, url string, data, target interface{}, reqOptions ...RequestOption) (int, error)
	GetWithQueryCustomHeader(c context.Context, url string, data, target interface{}, customHeader http.Header, reqOptions ...RequestOption) (int, error)
	GetWithoutEncodedQuery(c context.Context,
		url string, data, target interface{}, reqOptions ...RequestOption) (int, error)
	Do(ctx context.Context, request *http.Request, target interface{}, reqOptions ...RequestOption) (int, error)
	SendHTTPRequest(ctx context.Context, method string, path string, payload interface{}, outPut interface{}, reqOptions ...RequestOption) (int, error)
}

type client struct {
	client *http.Client
	opts   clientOptions
}

func NewClient(opts ...Option) Client {
	optsArg := getOptionsArg(opts)
	transport := NewTransport(optsArg)
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   optsArg.timeout,
	}
	c := &client{
		client: httpClient,
		opts:   optsArg,
	}
	return c
}

func getOptionsArg(opts []Option) clientOptions {
	// Init default options arg
	optsArgs := clientOptions{
		skipLog:         false,
		splitLogBody:    false,
		splitLogBodyLen: defaultLogBodyLength,
		timeout:         defaultTimeout,
	}

	for _, opt := range opts {
		opt.apply(&optsArgs)
	}
	return optsArgs
}

func (c *client) PostJSON(ctx context.Context,
	url string, data, target interface{}, reqOpts ...RequestOption) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(http.MethodPost).
		WithURL(url).
		WithHeaders(header).
		WithBody(MIMEJSON, data).
		Build()
	if err != nil {
		return
	}
	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) PostJSONWithCustomHeader(ctx context.Context,
	url string,
	data, target interface{},
	customHeader http.Header,
	method string,
	reqOpts ...RequestOption,
) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(method).
		WithURL(url).
		WithHeaders(header).
		WithBody(MIMEJSON, data).
		Build()
	if err != nil {
		return
	}
	req.Header = customHeader
	//decodeNumber := false
	//if len(reqOpts) > 0 {
	//	decodeNumber = reqOpts[0].DecodeNumber
	//}

	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) PostForm(ctx context.Context,
	url string, data, target interface{}, reqOpts ...RequestOption) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(http.MethodPost).
		WithURL(url).
		WithHeaders(header).
		WithBody(MIMEPOSTForm, data).
		Build()
	if err != nil {
		return
	}
	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) Get(ctx context.Context,
	url string, target interface{}, reqOpts ...RequestOption) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(http.MethodGet).
		WithURL(url).
		WithHeaders(header).
		Build()
	if err != nil {
		return
	}
	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) GetWithQuery(ctx context.Context,
	reqURL string, data, target interface{}, reqOpts ...RequestOption) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(http.MethodGet).
		WithURL(reqURL).
		WithHeaders(header).
		Build()
	if err != nil {
		return
	}

	if data != nil {
		v, err := query.Values(data)
		if err != nil {
			return 0, err
		}
		nonEncodedValue, _ := url.PathUnescape(v.Encode())
		req.URL.RawQuery = nonEncodedValue
	}
	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) GetWithoutEncodedQuery(ctx context.Context,
	reqURL string, data, target interface{}, reqOpts ...RequestOption) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(http.MethodGet).
		WithURL(reqURL).
		WithHeaders(header).
		Build()
	if err != nil {
		return
	}

	if data != nil {
		v, err := query.Values(data)
		if err != nil {
			return 0, err
		}
		nonEncodedValue, _ := url.QueryUnescape(v.Encode())
		req.URL.RawQuery = nonEncodedValue
	}
	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) Do(ctx context.Context, request *http.Request, target interface{}, reqOpts ...RequestOption) (int, error) {
	metricData := ginprometheus.GetPrometheusMetric()
	start := time.Now()
	if requestID := request.Header.Get(RequestIDHeader); requestID == "" {
		request.Header.Set(RequestIDHeader, getContextIDFromCtx(ctx))
	}

	decodeNumber := false
	if len(reqOpts) > 0 {
		decodeNumber = reqOpts[0].DecodeNumber
	}

	//handle StatusRequestTimeout
	var (
		bodyBytes []byte
		timer     = time.NewTimer(defaultTimeout)
	)
	url := request.URL.String()
	go func() {
		<-timer.C
		data := proccessBodyToResp(bodyBytes, url, request.Method)
		if strings.Contains(url, os.Getenv(constants.AWS_URL)) {
			message := fmt.Sprintf("your request has exceeded %v seconds", defaultTimeout.Seconds())
			c.sendGGchatErrorEvent(ctx,
				message,
				"The HTTP Gateway has timed out", data)
		}
	}()
	rsp, err := c.client.Do(request)
	if err != nil {
		timer.Stop()
		return 0, err
	}

	defer func() {
		_ = rsp.Body.Close()
		if strings.Contains(url, os.Getenv(constants.AWS_URL)) {
			if metricData != nil {
				go metricData.Latency(request, start, rsp.StatusCode, reqOpts[0].GroupPath)
			}

			//go metricData.Histogram(request, start)
			//go metricData.RequestTotal(request, rsp.StatusCode)
			//go metricData.RequestSize(request)
		}
	}()

	bodyBytes, err = io.ReadAll(rsp.Body)
	if err != nil {
		timer.Stop()
		return 0, err
	}
	timer.Stop()
	//handle StatusRequestTimeout
	if strings.Contains(string(bodyBytes), fmt.Sprintf(`"code":%v`, http.StatusRequestTimeout)) &&
		strings.Contains(url, os.Getenv(constants.AWS_URL)) {
		data := proccessBodyToResp(bodyBytes, url, request.Method)
		c.sendGGchatErrorEvent(ctx, "your request has timed out with body "+string(bodyBytes),
			"The HTTP Gateway has timed out", data)
	}

	if strings.Contains(string(bodyBytes), fmt.Sprintf(`"code":%v`, http.StatusTooManyRequests)) &&
		strings.Contains(url, os.Getenv(constants.AWS_URL)) {
		data := proccessBodyToResp(bodyBytes, url, request.Method)
		c.sendGGchatErrorEvent(ctx, "Too Many Requests "+string(bodyBytes),
			"Too Many Requests", data)
	}

	if len(bodyBytes) == 0 {
		return rsp.StatusCode, nil
	}

	if decodeNumber {
		d := json.NewDecoder(bytes.NewBuffer(bodyBytes))
		d.UseNumber()
		return rsp.StatusCode, d.Decode(target)
	}

	return rsp.StatusCode, json.Unmarshal(bodyBytes, target)
}

func (c *client) getRequestHeader(reqOpts ...RequestOption) map[string]string {
	if len(reqOpts) == 0 {
		return nil
	}
	reqOpt := reqOpts[0]
	header := reqOpt.Header
	if header == nil {
		header = make(map[string]string)
	}
	if reqOpt.GroupPath != "" {
		header[GroupPathHeader] = reqOpt.GroupPath
	}
	return header
}

func getContextIDFromCtx(ctx context.Context) string {
	if result, ok := ctx.Value(xcontext.KeyContextID.String()).(string); ok {
		return result
	}
	return ""
}

func (c *client) GetWithQueryCustomHeader(ctx context.Context,
	urlReq string, data, target interface{},
	customHeader http.Header, reqOpts ...RequestOption) (statusCode int, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err := NewRequestBuilderWithCtx(ctx).
		WithMethod(http.MethodGet).
		WithURL(urlReq).
		WithHeaders(header).
		Build()
	if err != nil {
		return
	}

	if data != nil {
		v, err := query.Values(data)
		if err != nil {
			return 0, err
		}
		nonEncodedValue, _ := url.PathUnescape(v.Encode())
		req.URL.RawQuery = nonEncodedValue
		// req.URL.RawQuery = v.Encode()
	}
	fmt.Println("111111111111111111111111111122221", data)
	fmt.Println("11111111111111111111111111111", req.URL)
	req.Header = customHeader
	//decodeNumber := false
	//if len(reqOpts) > 0 {
	//	decodeNumber = reqOpts[0].DecodeNumber
	//}

	return c.Do(ctx, req, target, reqOpts...)
}

func (c *client) SendHTTPRequest(
	ctx context.Context,
	method string,
	path string,
	payload interface{},
	outPut interface{},
	reqOptions ...RequestOption,
) (status int, err error) {
	req, err := c.newRequest(ctx, method, path, payload, reqOptions)
	if err != nil {
		return -1, fmt.Errorf("failed to create %s request: %w", method, err)
	}

	status, err = c.doRequest(ctx, req, outPut)
	if err != nil {
		return
	}

	return
}

/*Internal implementation*/
func (c *client) newRequest(
	ctx context.Context,
	method,
	path string,
	payload interface{},
	reqOpts []RequestOption,
) (req *http.Request, err error) {
	header := c.getRequestHeader(reqOpts...)
	req, err = NewRequestBuilderWithCtx(ctx).
		WithMethod(method).
		WithBody(header[contentTypeField], payload).
		// WithURL(fmt.Sprintf("%s/%s", strings.TrimRight(c.Ops.URL, "/"), path)).
		WithURL(path).
		WithHeaders(header).
		Build()
	if err != nil {
		return
	}
	if reqOpts[0].HeaderCustom != nil {
		req.Header = reqOpts[0].HeaderCustom
	}

	return req, nil
}

type DataError struct {
	Mess string `json:"mess"`
	Code int    `json:"code"`
	Data struct {
		Code    string `json:"code"`
		Details string `json:"details"`
	} `json:"data"`
	Y4AMess string `json:"y4aMess"`
	Y4ACode int    `json:"y4aCode"`
}

func (h *client) doRequest(
	ctx context.Context,
	r *http.Request,
	outPut interface{},
) (status int, err error) {
	if requestID := r.Header.Get(RequestIDHeader); requestID == "" {
		r.Header.Set(RequestIDHeader, getContextIDFromCtx(ctx))
	}

	apmClient := apmhttp.WrapClient(h.client)
	resp, err := ctxhttp.Do(ctx, apmClient, r)
	if err != nil {
		logger.NewLogger().WithKeyword(ctx, "MAKE_REQUEST_ERROR").
			WithFields(
				logrus.Fields{
					"URL":    r.URL.String(),
					"Method": r.Method,
				}).
			WithStatusCode(resp.StatusCode).
			WithError(err).
			Error()
		return -1, fmt.Errorf("failed to make request: %w", err)
	}

	if resp == nil {
		return
	}

	// return first if not need output return
	if outPut == nil {
		return
	}

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err := dec.Decode(outPut); err != nil {
		logger.NewLogger().WithKeyword(ctx, "PARSE_RESPONSE_BODY_ERROR").
			WithFields(
				logrus.Fields{
					"URL":    r.URL.String(),
					"Method": r.Method,
					"Output": buf.String(),
				}).
			WithStatusCode(resp.StatusCode).
			WithError(err).
			Error()
		return resp.StatusCode, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.NewLogger().WithKeyword(ctx, "DO_HTTP_REQUEST_ERROR").
			WithFields(
				logrus.Fields{"Status": resp.Status,
					"PostForm": r.PostForm,
					"Form":     r.Form,
					"Header":   r.Header,
					"URL":      r.URL.String(),
					"Method":   r.Method,
					"Output":   outPut,
				}).
			WithStatusCode(resp.StatusCode).
			WithError(err).
			Error()

		switch resp.StatusCode {
		case http.StatusInternalServerError:
			return
		default:
			return status, errors.New(resp.Status)
		}
	}

	defer resp.Body.Close() // nolint: errcheck

	return http.StatusOK, nil
}

func (c *client) sendGGchatErrorEvent(ctx context.Context,
	errMessage string, topic string, data dataResponse) {

	trackID := tracking.GetTrackIDFromContext(ctx)
	env := os.Getenv(constants.EnvironmentsEnv)
	path := GetEnviromentHttpTimeout(env)
	fromSource, _ := xcontext.GetFromSource(ctx)
	if env == "" {
		return
	}
	message := fmt.Sprintf("```Log data started at: %v \n \ntopic: %v \nenv: %v \nfrom_source: %v \nerror: %v \n"+
		"url: %v \nmethod: %v \nrequest_id: %v \ntrack_id: %v```",
		timeutils.NowInGMT07().Format(timeutils.CustomDateTimeVN), topic, env, fromSource, errMessage,
		data.URL, data.Method, data.RequestId, trackID)

	go c.sendGGChatCustom(ctx,
		message,
		path)
}

func (c *client) sendGGChatCustom(
	ctx context.Context,
	message string,
	path string,

) (err error) {
	messageData := struct {
		Text string `json:"text"`
	}{
		Text: message,
	}
	var returnValue ggResponse
	xopt := RequestOption{GroupPath: "/chat.googleapis.com/v1/spaces"}
	if _, err = c.PostJSON(ctx, path, &messageData, &returnValue, xopt); err != nil {
		apm.CaptureError(context.Background(), err).Send()
		return
	}

	return
}

func proccessBodyToResp(body []byte, url string, method string) (dataResp dataResponse) {
	err := json.Unmarshal(body, &dataResp)
	if err != nil {
		dataResp.URL = url
		dataResp.Method = method
		return
	}
	dataResp.URL = url
	dataResp.Method = method
	prefix := "requestId: "

	startIndex := strings.Index(dataResp.RequestId, prefix)
	if startIndex == -1 {
		return
	}
	dataResp.RequestId = dataResp.RequestId[startIndex+len(prefix):]
	return dataResp
}
