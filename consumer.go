package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func runConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Consumer shutting down")
		wg.Done()
	}()

	tp, err := makeTraceProvider("consumer")
	if err != nil {
		panic(err)
	}

	cg, _, err := getKafkaConsumerGroup()
	if err != nil {
		panic(err)
	}
	defer cg.Close()

	log.Println("Consumer started")
	handler := otelsarama.WrapConsumerGroupHandler(
		&topicHandler{
			tp: tp,
		},
		otelsarama.WithTracerProvider(tp),
	)
	for {
		if err := cg.Consume(ctx, []string{kafkaTopic}, handler); err != nil {
			log.Printf("[ERR] consume failed. ERR: %+v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

type topicHandler struct {
	tp *tracesdk.TracerProvider
}

func (h *topicHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *topicHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *topicHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		for _, h := range msg.Headers {
			log.Printf("[consumer] header `%s`: %s", string(h.Key), string(h.Value))
		}
		log.Printf("[consumer] key: `%s`", string(msg.Key))
		log.Printf("[consumer] message: %s", string(msg.Value))

		ctx := session.Context()
		traceCtx := otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewConsumerMessageCarrier(msg))
		doTrace(
			traceCtx,
			h.tp,
			"consume-task",
			nil,
			func(ctx context.Context, span trace.Span) error {
				time.Sleep(50 * time.Millisecond)
				return nil
			},
		)
		session.MarkMessage(msg, "")
	}
	return nil
}
