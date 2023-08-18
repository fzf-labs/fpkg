package tracing

import (
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// NewJaegerExporter 创建一个jaeger导出器
func NewJaegerExporter(endpoint string) (traceSdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
}

// NewZipkinExporter 创建一个zipkin导出器
func NewZipkinExporter(endpoint string) (traceSdk.SpanExporter, error) {
	return zipkin.New(endpoint)
}

// NewTracerExporter 创建一个导出器，支持：jaeger和zipkin
func NewTracerExporter(exporterName, endpoint string) (traceSdk.SpanExporter, error) {
	if exporterName == "" {
		exporterName = "jaeger"
	}
	switch exporterName {
	case "jaeger":
		return NewJaegerExporter(endpoint)
	case "zipkin":
		return NewZipkinExporter(endpoint)
	default:
		return nil, errors.New("exporter type not support")
	}
}

type Config struct {
	ServiceName  string
	ExporterName string
	Endpoint     string
	Sampler      float64
	Version      string
	InstanceID   string
	Env          string
}

// NewTracerProvider 创建一个链路追踪器
func NewTracerProvider(cfg *Config) *traceSdk.TracerProvider {
	if cfg == nil {
		return nil
	}
	if cfg.Sampler == 0 {
		cfg.Sampler = 1.0
	}
	if cfg.Env == "" {
		cfg.Env = "dev"
	}
	opts := []traceSdk.TracerProviderOption{
		traceSdk.WithSampler(traceSdk.ParentBased(traceSdk.TraceIDRatioBased(cfg.Sampler))),
		traceSdk.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(cfg.ServiceName),
			semConv.ServiceVersionKey.String(cfg.Version),
			semConv.ServiceInstanceIDKey.String(cfg.InstanceID),
			attribute.String("env", cfg.Env),
		)),
	}
	if len(cfg.Endpoint) > 0 {
		exp, err := NewTracerExporter(cfg.ExporterName, cfg.Endpoint)
		if err != nil {
			panic(err)
		}
		opts = append(opts, traceSdk.WithBatcher(exp))
	}
	return traceSdk.NewTracerProvider(opts...)
}
