package notification

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"
	kafkaClient "yes4all/ads-noti-api/pkg/queue/kafka"
	"yes4all/ads-noti-api/pkg/utils/constants"
	"yes4all/ads-noti-api/services/ads-noti/repository"
	"yes4all/ads-noti-api/services/ads-noti/usecase"
)

var InitKafkaConn = fx.Invoke(
	initKafkaConn,
)

var SubsribeToTopic = fx.Invoke(
	kafkaSubscribeToTopic,
)

var SubsribeToTopicDLQ = fx.Invoke(
	SubscribeAndRepushNotification,
)

func initKafkaConn() {
	maxRetries := constants.MaxRetryConnectKafkaTimes                     // Số lần thử tối đa
	retryInterval := constants.MaxRetryConnectKafkaDuration * time.Second // Thời gian chờ giữa các lần thử lại

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancel()

	kafkaConfig := config.KafkaConfig()
	loggerCustom := logger.NewLogger()

	var (
		conn *kafka.Conn
		err  error
	)

	for retry := 0; retry < maxRetries; retry++ {
		conn, err = kafkaClient.ConnectKafkaBrokers(ctx, kafkaConfig)
		if err == nil {
			// Kết nối thành công, thoát khỏi vòng lặp
			break
		}

		loggerCustom.WithKeyword(ctx, "connectKafkaBrokers got error").WithError(err).Warningf("Retrying in %s...", retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		// Không thể kết nối sau số lần thử tối đa
		loggerCustom.WithKeyword(ctx, "connectKafkaBrokers got error").WithError(err).Fatal("Max retries reached. Could not connect to Kafka.")
	}

	if kafkaConfig.InitTopics {
		kafkaClient.InitKafkaTopics(ctx, conn, kafkaTopicConfig(kafkaConfig)...)
	}
}

// init kafka topic dlq
func kafkaTopicConfig(
	kafkaConfig config.KafkaCfg,
) (topics []kafka.TopicConfig) {
	//NumPartitions:     1, để chạy tuần tự
	//ReplicationFactor: 1,
	topics = append(topics,
		kafka.TopicConfig{
			Topic:             kafkaConfig.TopicNotification,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		kafka.TopicConfig{
			Topic:             kafkaConfig.TopicNotificationDLQ,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	return
}

func kafkaSubscribeToTopic() {
	kafkaConfig := config.KafkaConfig()
	log := logger.NewLogger()
	processMessage := NewProductMessageProcessor(
		usecase.NewWebsocketUsecase(repository.NewNotificationRepo()))

	cg := kafkaClient.NewConsumerGroup(kafkaConfig.GetBrokers(), kafkaConfig.GroupID, log)
	// start comsume message
	go cg.ConsumeTopic(kafkaConfig.GetTopicsNotification(), 1, processMessage.ProcessMessages)
}
