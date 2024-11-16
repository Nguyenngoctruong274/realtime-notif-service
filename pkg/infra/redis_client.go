package infra

import (
	"context"
	"log"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/go-redis/redis/v8"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
)

type RedisClient struct {
	Client *redis.Client
}

var redisClientSingleton *RedisClient

func InitRedisClient() {
	var ctx = context.Background()

	redisURL := config.RedisConfig().RedisURL
	entry := logger.NewLogger().WithKeyword(ctx, "InitRedisClient").
		WithField("redisURL", redisURL)
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set,
	})

	// Check heath
	_, err := client.Ping(ctx).Result()
	if err != nil {
		entry.WithError(err).Error()
		return
	}

	client.AddHook(apmgoredis.NewHook())
	info, err := client.Info(ctx).Result()
	entry.WithField("info", info)
	if err != nil {
		entry.WithError(err).Error()
		return
	}
	entry.Info()

	redisClientSingleton = &RedisClient{Client: client}
}

func GetRedisClient() *redis.Client {
	if redisClientSingleton == nil {
		log.Fatal("failed to create new redis client")
	}

	return redisClientSingleton.Client
}

func CloseRedis() {
	redisClientSingleton.Client.Close()
}
