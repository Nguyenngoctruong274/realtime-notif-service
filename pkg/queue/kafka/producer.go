package kafka

import (
	"context"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"

	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

var Module = fx.Invoke(initKafkaProducer)

type KafkaWriter struct {
	KafkaWriter *kafka.Writer
}

var kafkaWriterSigleton *KafkaWriter

func initKafkaProducer() {
	log := logger.NewLogger()
	kafkaCfg := config.KafkaConfig()
	brokers := kafkaCfg.GetBrokers()
	kafkaWriter := NewWriter(brokers, kafka.LoggerFunc(log.WithInput("NewWriter").Errorf))
	if kafkaWriter == nil {
		log.Fatal("Failed to initialize kafka writer.")
	}
	kafkaWriterSigleton = &KafkaWriter{kafkaWriter}
}

func PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	if kafkaWriterSigleton == nil {
		initKafkaProducer()
	}
	logger.NewLogger().WithKeyword(ctx, "KAFKA_PUBLISH_MESSAGE").WithOutput(msgs).Info()
	return kafkaWriterSigleton.KafkaWriter.WriteMessages(ctx, msgs...)
}

func Close() error {
	return kafkaWriterSigleton.KafkaWriter.Close()
}
