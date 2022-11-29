package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(e *echo.Echo, tp *tracesdk.TracerProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.WrapHandler(otelhttp.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				echoCtx := e.NewContext(r, w)
				next(echoCtx)
			}),
			"http.handler",
			otelhttp.WithTracerProvider(tp),
		))
	}
}

func LogTraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := trace.SpanFromContext(c.Request().Context())
			if span.IsRecording() {
				traceId := span.SpanContext().TraceID().String()
				traceUrl := fmt.Sprintf("http://localhost:16686/trace/%s", traceId)
				c.Response().Header().Set("X-Trace-Id", traceId)
				c.Response().Header().Set("X-Trace-Url", traceUrl)
			}
			return next(c)
		}
	}
}
