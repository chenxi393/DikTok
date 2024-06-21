package otel

import (
	"context"
	"log"
	"os"

	"diktok/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// 这个可选使用 其实我认为使用比较好
// 业务代码不要引入别的乱起八糟的包 抽象出来比较好
var Tracer trace.Tracer

func newStdoutExporter(serviceName string) *stdout.Exporter {
	filePath := "./" + serviceName + ".log"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln("无法打开文件:", err)
	}
	// defer file.Close()
	exporter, err := stdout.New(stdout.WithPrettyPrint(), stdout.WithWriter(file))
	if err != nil {
		log.Fatal(err)
	}
	return exporter
}

func newExporterGRPC() *otlptrace.Exporter {
	// Your preferred exporter: console, jaeger, zipkin, OTLP, etc.
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			// 很奇怪 grpc 不能使用 sock5 代理 debug很久才发现 走的代理
			otlptracegrpc.WithEndpoint(config.System.OtelColletcor.Host+":"+config.System.OtelColletcor.Port),
		),
	)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	return exporter
}

func newExporterHTTP() *otlptrace.Exporter {
	// Your preferred exporter: console, jaeger, zipkin, OTLP, etc.
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			// 使用http 这里正常的 不使用https
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(config.System.OtelColletcor.Host+":"+config.System.OtelColletcor.Port),
		),
	)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	return exporter
}

// Create a new tracer provider with a batch span processor and the given exporter.
func initTracerProvider(url, serviceName string) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(newExporterHTTP()),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				url,
				semconv.ServiceNameKey.String(serviceName),
			)),
	)
	// 这里设置为全局的
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Finally, set the tracer that can be used for this package.
	Tracer = tp.Tracer(serviceName)
	return tp
}
