package kafka

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func ConnectKafkaBrokers(ctx context.Context, cfg config.KafkaCfg) (conn *kafka.Conn, err error) {
	loggerCustom := logger.NewLogger()
	conn, err = NewKafkaConn(ctx, &Config{
		Brokers:    cfg.GetBrokers(),
		GroupID:    cfg.GroupID,
		InitTopics: cfg.InitTopics,
	}) // TODO:
	if err != nil {
		return nil, errors.Wrap(err, "kafka.NewKafkaCon")
	}

	brokers, err := conn.Brokers()
	if err != nil {
		return nil, errors.Wrap(err, "kafkaConn.Brokers")
	}

	loggerCustom.WithKeyword(
		ctx,
		fmt.Sprintf("kafka connected to brokers: %+v", brokers)).
		Info()

	return
}

func InitKafkaTopics(ctx context.Context, kafkaConn *kafka.Conn, topics ...kafka.TopicConfig) {
	loggerCustom := logger.NewLogger()

	controller, err := kafkaConn.Controller()
	if err != nil {
		loggerCustom.Warnf("kafkaConn.Controller err: %v", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	loggerCustom.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		loggerCustom.Warnf("initKafkaTopics.DialContext: %v", err)
		return
	}
	defer conn.Close() //nolint: errcheck

	loggerCustom.Infof("established new kafka controller connection: %s", controllerURI)

	if err := conn.CreateTopics(
		topics...,
	); err != nil {
		loggerCustom.Warnf("kafkaConn.CreateTopics: %v", err)
		return
	}

	loggerCustom.Infof("kafka topics created or already exists: %+v", topics)
}

func LogProcessMessage(ctx context.Context, m kafka.Message, workerID int) {
	logger.NewLogger().
		WithKeyword(ctx, "Kafka log process").
		WithField("Topic: ", m.Topic).
		WithField("Partition: ", m.Partition).
		WithField("WorkerID: ", workerID).
		WithField("Offset: ", m.Offset).
		WithTime(m.Time).
		WithField("ValueMessage : ", string(m.Value)).
		Info()
}
