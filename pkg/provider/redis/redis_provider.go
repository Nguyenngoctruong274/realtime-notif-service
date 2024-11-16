package provider

import (
	"context"
	"time"
	"yes4all/ads-noti-api/pkg/infra"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/utils"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
)

type redisProvider struct {
	client *redis.Client
	logger *logger.StandardLogger
}

func NewRedisProvider(redisClient infra.RedisClient) IRedisProvider {
	return &redisProvider{
		client: redisClient.Client,
		logger: logger.NewLogger(),
	}
}

type IRedisProvider interface {
	HSetData(ctx context.Context, key string, data interface{}, durationTime time.Duration) error
	HGetAllData(ctx context.Context, key string, data interface{}) (bool, error)
	Remove(ctx context.Context, keys ...string) error
}

func (r *redisProvider) HSetData(ctx context.Context, key string, data interface{}, durationTime time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisProvider.HSetData")
	span.SetTag("param.key", key)
	span.SetTag("param.durationTime", durationTime)
	span.SetTag("param.data", data)

	mapData := utils.Inspect(data)
	_, err := r.client.HSet(ctx, key, mapData).Result()
	if err != nil {
		r.logger.WithKeyword(ctx, "HSetData_HSet").
			WithInput(map[string]interface{}{"key": key, "mapData": mapData}).
			WithError(err).
			Error()

		return err
	}
	if durationTime > 0 {
		_, err = r.client.Expire(ctx, key, durationTime).Result()
		if err != nil {
			r.logger.WithKeyword(ctx, "HSetData_Expire").
				WithInput(map[string]interface{}{"key": key, "durationTime": durationTime, "mapData": mapData}).
				WithError(err).
				Error()

			return err
		}
	}

	r.logger.WithKeyword(ctx, "HSetData").
		WithInput(map[string]interface{}{"key": key, "mapData": mapData}).
		Info()

	return nil
}

func (r *redisProvider) HGetAllData(ctx context.Context, key string, data interface{}) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisProvider.HGetAllData")
	span.SetTag("param.key", key)
	span.SetTag("param.data", data)

	redisData, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		r.logger.WithKeyword(ctx, "HGetAllData").
			WithInput(map[string]interface{}{"key": key}).
			WithError(err).
			Error()

		return false, err
	}
	r.logger.WithKeyword(ctx, "HGetAllData").
		WithInput(map[string]interface{}{"key": key, "redisData": redisData}).
		Info()
	if len(redisData) < 1 {
		return true, nil
	}

	cfg := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           data,
		WeaklyTypedInput: true,
		ZeroFields:       true,
		ErrorUnused:      false,
		TagName:          `json`,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			utils.CustomDecodeHook()),
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err = decoder.Decode(redisData)
	if err != nil {
		r.logger.WithKeyword(ctx, "HGetAllData_Decode").
			WithInput(map[string]interface{}{"key": key, "response": redisData}).
			WithError(err).
			Error()

		return false, err
	}

	return false, nil
}

func (r *redisProvider) Remove(ctx context.Context, keys ...string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisProvider.Remove")
	span.SetTag("param.keys", keys)

	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.logger.WithKeyword(ctx, "Remove").
			WithInput(keys).
			WithError(err).
			Error()

		return err
	}
	r.logger.WithKeyword(ctx, "HSetData").
		WithInput(map[string]interface{}{"keys": keys}).
		Info()

	return nil
}
