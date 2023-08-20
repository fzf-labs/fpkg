package httputil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/tracing"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	client := NewClient()
	tracerProvider := tracing.NewTracerProvider(&tracing.Config{
		ServiceName:  "http",
		ExporterName: "jaeger",
		Endpoint:     "http://127.0.0.1:14268/api/traces",
		Sampler:      1.0,
		Version:      "123",
		InstanceID:   "123456",
		Env:          "dev",
	})
	client.SetTracer(tracerProvider.Tracer("baidu"))
	response, err := client.R().SetContext(ctx).Post("http://www.baidu.com")
	fmt.Println(response.Status)
	fmt.Println(err)
	time.Sleep(time.Second * 10)
	assert.Equal(t, nil, err)
}
