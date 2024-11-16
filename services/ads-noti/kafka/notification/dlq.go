package notification

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"
	kafkaClient "yes4all/ads-noti-api/pkg/queue/kafka"
	"yes4all/ads-noti-api/services/ads-noti/repository"
	"yes4all/ads-noti-api/services/ads-noti/usecase"
)

type retryMessageProcessor struct {
	socketUsecase usecase.WebsocketUsecase
	kafkConfig    config.KafkaCfg
}

type IretryMessageProcessor interface {
	ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
}

func NewRetryMessageProcessor(
	socketUsecase usecase.WebsocketUsecase,
) IretryMessageProcessor {
	return &retryMessageProcessor{}
}

func (p *retryMessageProcessor) ProcessMessages(
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
		case config.KafkaConfig().TopicNotificationDLQ:
			fmt.Println(string(m.Value))
			p.ProcessRetryDLQ(ctx, r, m)
		}
	}
}

func (s *retryMessageProcessor) ProcessRetryDLQ(
	ctx context.Context,
	r *kafka.Reader,
	m kafka.Message,
) {
	kafkaConfig := config.KafkaConfig()
	// do something stuff
	var retryOptions = []retry.Option{
		retry.Attempts(kafkaConfig.KafkaRetryAttempts),
		retry.Delay(time.Duration(kafkaConfig.KafkaRetryDelay) * time.Second),
		retry.DelayType(retry.BackOffDelay)}

	message := kafka.Message{
		Topic: kafkaConfig.TopicNotification,
		Value: m.Value,
		Key:   []byte(s.kafkConfig.DLQMessageKey),
	}

	if err := retry.Do(func() error {
		return kafkaClient.PublishMessage(ctx, message)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.commitMessage(ctx, r, m)
		return
	}

	s.commitMessage(ctx, r, m)
}

func SubscribeAndRepushNotification() {
	//ctx, cancel := signal.NotifyContext(
	//	context.Background(),
	//	os.Interrupt,
	//	syscall.SIGTERM,
	//	syscall.SIGINT,
	//)
	//defer cancel()
	kafkaConfig := config.KafkaConfig()
	log := logger.NewLogger()
	processMessage := NewRetryMessageProcessor(
		usecase.NewWebsocketUsecase(repository.NewNotificationRepo()))

	cg := kafkaClient.NewConsumerGroup(kafkaConfig.GetBrokers(), "canh_test1", log)
	// start comsume message
	go cg.ConsumeTopic(
		[]string{kafkaConfig.TopicNotificationDLQ},
		1,
		processMessage.ProcessMessages,
	)

	//<-ctx.Done()
}

func (s *retryMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	loggerCustom := logger.NewLogger()
	if err := r.CommitMessages(ctx, m); err != nil {
		loggerCustom.Warnf("commitMessage: %v", err)
	}
}

func (s *retryMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	loggerCustom := logger.NewLogger()
	// TODO: check error

	err := s.publishToDLQ(ctx, m)
	if err != nil {
		loggerCustom.WithKeyword(ctx, "error publishing to dlq").Error(err)
	}

	if err := r.CommitMessages(ctx, m); err != nil {
		loggerCustom.Warnf("commitErrMessage: %v", err)
	}
}

func (s *retryMessageProcessor) publishToDLQ(
	ctx context.Context,
	m kafka.Message,
) error {
	err := kafkaClient.PublishMessage(ctx, kafka.Message{
		Topic: s.kafkConfig.TopicNotificationDLQ,
		Key:   []byte("DLQ"),
		Value: m.Value,
	})
	if err != nil {
		logger.NewLogger().WithKeyword(ctx, "publishToDLQ error").Error()
		return err
	}

	return err
}
