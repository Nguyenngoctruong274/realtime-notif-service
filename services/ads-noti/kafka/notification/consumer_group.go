package notification

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"sync"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"
	kafkaClient "yes4all/ads-noti-api/pkg/queue/kafka"
	"yes4all/ads-noti-api/services/ads-noti/usecase"
)

type productMessageProcessor struct {
	socketUsecase usecase.WebsocketUsecase
	kafkConfig    config.KafkaCfg
}

type IprocessMessage interface {
	ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
}

func NewProductMessageProcessor(
	socketUsecase usecase.WebsocketUsecase,
) IprocessMessage {
	return &productMessageProcessor{
		socketUsecase: socketUsecase,
		kafkConfig:    config.KafkaConfig(),
	}
}

func (p *productMessageProcessor) ProcessMessages(
	ctx context.Context,
	r *kafka.Reader,
	wg *sync.WaitGroup,
	workerID int,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		kafkaClient.LogProcessMessage(ctx, m, workerID)
		switch m.Topic {
		case config.KafkaConfig().TopicNotification:
			p.ProcessNotification(ctx, r, m)
		}
	}
}

func (s *productMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	loggerCustom := logger.NewLogger()
	if err := r.CommitMessages(ctx, m); err != nil {
		loggerCustom.Warnf("commitMessage: %v", err)
	}
}

func (s *productMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	loggerCustom := logger.NewLogger()
	// TODO: check error

	err := s.publishToDLQ(ctx, m)
	if err != nil {
		loggerCustom.WithKeyword(ctx, "error publishing to dlq").Error(err)
	}

	if err = r.CommitMessages(ctx, m); err != nil {
		loggerCustom.Warnf("commitErrMessage: %v", err)
	}
}

func (s *productMessageProcessor) publishToDLQ(
	ctx context.Context,
	m kafka.Message,
) error {
	smg := kafka.Message{
		Topic: s.kafkConfig.TopicNotificationDLQ,
		Value: m.Value,
		Key:   []byte(s.kafkConfig.DLQMessageKey),
	}
	err := kafkaClient.PublishMessage(ctx, smg)
	if err != nil {
		logger.NewLogger().WithKeyword(ctx, "publishToDLQ error").Error()
		return err
	}
	kafkaClient.LogProcessMessage(ctx, smg, 0)

	return err
}
