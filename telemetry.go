package main

import (
	"context"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

func makeTraceProvider(svcName string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(jaegerEndpoint)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		)),
	)
	return tp, nil
}

func doTrace(
	ctx context.Context,
	tp *tracesdk.TracerProvider,
	name string,
	spanFn func(trace.Span),
	fn func(context.Context) error,
) error {
	spanCtx, span := tp.Tracer("vanhtuan0409.tracing").Start(ctx, name)
	defer span.End()

	if spanFn != nil {
		spanFn(span)
	}
	err := fn(spanCtx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
