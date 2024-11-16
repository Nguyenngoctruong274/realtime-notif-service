package notification

import (
	"context"
	"encoding/json"
	"time"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/services/ads-noti/model/request"

	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"
)

func (s *productMessageProcessor) ProcessNotification(
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

	var req request.MarkAsReadRequest
	json.Unmarshal(m.Value, &req)

	// retry 3 times
	if err := retry.Do(func() error {
		return s.socketUsecase.SubData(ctx, req)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		// if message is retry from DLQ then commit and return
		if m.Key != nil {
			s.commitMessage(ctx, r, m)
			return
		}
		// else push to DLQ and then try again
		s.commitErrMessage(ctx, r, m)
		return
	}
	// just commit here if successfully
	s.commitMessage(ctx, r, m)
}
