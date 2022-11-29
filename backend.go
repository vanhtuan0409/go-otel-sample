package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func runBackend(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Backend shutting down")
		wg.Done()
	}()

	tp, err := makeTraceProvider("backend")
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	go func() {
		<-ctx.Done()
		e.Shutdown(context.Background())
	}()

	e.Use(TracingMiddleware(e, tp))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "backend")
	})
	e.GET("/echo", echoHandler(tp))

	log.Println("Backend starting up")
	addr := fmt.Sprintf(":%d", backendPort)
	e.Start(addr)
}

func echoHandler(tp *tracesdk.TracerProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		resp := struct {
			Request string
			Header  http.Header
		}{
			Request: req.URL.String(),
			Header:  req.Header,
		}

		doTrace(
			req.Context(),
			tp,
			"heavy-task",
			nil,
			func(ctx context.Context, span trace.Span) error {
				time.Sleep(50 * time.Millisecond)
				return nil
			},
		)
		return c.JSON(http.StatusOK, &resp)
	}
}
