package httputil

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	//ctx := context.Background()
	//traceProvider, err := jaegerTraceProvider("github-query", "test", "http://127.0.0.1:14268/api/traces")
	//if err != nil {
	//	return
	//}
	//defer traceProvider.Shutdown(ctx)
	//otel.SetTracerProvider(traceProvider)
	//client := NewClient()
	//client.SetR(client, otel.Tracer("github"))
	//response, err := client.R().SetContext(ctx).Get("https://www.baidu.com/")
	//if err != nil {
	//	return
	//}
	//fmt.Println(response)
}
