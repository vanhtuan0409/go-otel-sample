package main

import (
	"github.com/Shopify/sarama"
)

func getKafkaConsumerGroup() (sarama.ConsumerGroup, *sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cg, err := sarama.NewConsumerGroup(kafkaBrokers, "otel-sample", cfg)
	return cg, cfg, err
}

func getKafkaProducer() (sarama.SyncProducer, *sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Return.Successes = true
	cfg.Producer.Compression = sarama.CompressionGZIP
	p, err := sarama.NewSyncProducer(kafkaBrokers, cfg)
	return p, cfg, err
}
