package otel

import (
	"context"
	"douyin/config"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initMeterProvider(url, serviceName string) *sdkmetric.MeterProvider {
	ctx := context.Background()
	exporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(config.System.OtelColletcor.Host+":"+config.System.OtelColletcor.Port),
	)
	if err != nil {
		log.Fatalf("new otlp metric exporter failed: %v", err)
	}
	// Ensure default SDK resources and the required service name are set.
	r := sdkmetric.WithResource(
		resource.NewWithAttributes(
			url,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		r,
		sdkmetric.WithView(sdkmetric.NewView(
			sdkmetric.Instrument{Scope: instrumentation.Scope{Name: "go.opentelemetry.io/contrib/google.golang.org/grpc/otelgrpc"}},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationDrop{}},
		)),
	)
	otel.SetMeterProvider(mp)
	return mp
}
