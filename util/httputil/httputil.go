package httputil

import (
	"github.com/imroc/req/v3"
)

type Client struct {
	*req.Client
}

func NewClient() *Client {
	client := req.C().EnableTraceAll().EnableDumpEachRequest().OnAfterResponse(ResponseMiddleware)
	return &Client{
		client,
	}
}

// SetDebug enable debug if set to true, disable debug if set to false.
func (c *Client) SetDebug(enable bool) *Client {
	if enable {
		c.EnableDebugLog()
		c.EnableDumpAll()
	} else {
		c.DisableDebugLog()
		c.DisableDumpAll()
	}
	return c
}

// SetPrometheus 设置Prometheus
func (c *Client) SetPrometheus() *Client {
	c.OnAfterResponse(Prometheus())
	return c
}

// SetTracer 设置Prometheus
func (c *Client) SetTracer() *Client {
	jaegerTraceProvider()
	return c
}
