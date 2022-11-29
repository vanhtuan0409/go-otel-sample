package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{}, // must add this to probagate parent context
	))
}

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
	startOptions []trace.SpanStartOption,
	fn func(context.Context, trace.Span) error,
) error {
	tracer := tp.Tracer("vanhtuan0409.tracing")
	spanCtx, span := tracer.Start(ctx, name, startOptions...)
	defer span.End()

	err := fn(spanCtx, span)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
