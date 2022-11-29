package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
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

	p, cfg, err := getKafkaProducer()
	if err != nil {
		panic(err)
	}
	p = otelsarama.WrapSyncProducer(cfg, p, otelsarama.WithTracerProvider(tp))
	defer p.Close()

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
	e.GET("/echo", echoHandler(tp, p))

	log.Println("Backend starting up")
	addr := fmt.Sprintf(":%d", backendPort)
	e.Start(addr)
}

func echoHandler(tp *tracesdk.TracerProvider, p sarama.SyncProducer) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		doTrace(
			req.Context(),
			tp,
			"backend-task",
			nil,
			func(ctx context.Context, s trace.Span) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		)

		msg := &sarama.ProducerMessage{
			Topic: kafkaTopic,
			Value: sarama.StringEncoder("hello world!!!"),
		}
		otel.GetTextMapPropagator().Inject(req.Context(), otelsarama.NewProducerMessageCarrier(msg))
		p.SendMessage(msg)

		resp := struct {
			Request string
			Header  http.Header
		}{
			Request: req.URL.String(),
			Header:  req.Header,
		}
		return c.JSON(http.StatusOK, &resp)
	}
}
