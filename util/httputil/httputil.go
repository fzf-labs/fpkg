package httputil

import "github.com/imroc/req/v3"

func NewClient() {
	req.C().EnableTraceAll()
}
