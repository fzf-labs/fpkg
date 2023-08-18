package httputil

import (
	"strconv"
	"time"

	"github.com/imroc/req/v3"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	SendHTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "send_http_requests_total",
			Help: "Number of the http requests sent since the server started",
		},
		[]string{"method", "host", "path", "code"},
	)
	SendHTTPRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "send_http_requests_duration_seconds",
			Help:    "Duration in seconds to send http requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "host", "path", "code"},
	)
)

//nolint:gochecknoinits
func init() {
	prometheus.MustRegister(SendHTTPRequests, SendHTTPRequestsDuration)
}

// Prometheus 中间件统一记录 Prometheus 指标
func Prometheus() func(c *req.Client, resp *req.Response) error {
	return func(c *req.Client, resp *req.Response) error {
		request := resp.Request
		code := ""
		if resp.Response != nil {
			code = strconv.Itoa(resp.Response.StatusCode)
		}
		SendHTTPRequests.WithLabelValues(
			request.Method, request.URL.Host, request.URL.Path, code,
		).Inc()

		duration := float64(resp.TotalTime()) / float64(time.Second)
		SendHTTPRequestsDuration.WithLabelValues(
			request.Method, request.URL.Host, request.URL.Path, code,
		).Observe(duration)
		return nil
	}
}

// SetPrometheus 设置Prometheus
func (c *Client) SetPrometheus() *Client {
	c.OnAfterResponse(Prometheus())
	return c
}
