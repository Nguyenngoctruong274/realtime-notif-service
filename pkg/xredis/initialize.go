package xredis

import (
	"strings"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/go-redis/redis"
)

type Config struct {
	SingleMode bool
	URL        string
	Password   string
}

func GetRedisClient(redisURL, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: password,
	})

	_, err := client.Ping().Result()
	return client, err
}

func GetRedisClusterClient(redisURLs string, password string) (*redis.ClusterClient, error) {
	addrs := strings.Split(redisURLs, ",")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})

	_, err := client.Ping().Result()
	return client, err
}

type Closable interface {
	Close() error
}

func CloseRedis(c Closable) func() {
	return func() {
		err := c.Close()
		if err != nil {
			logger.NewLogger().Errorf("Error: failed to close redis, error: %v", err)
		}
	}
}

func Init(config Config) (Cmdable, func(), error) {
	if config.SingleMode {
		redisClient, err := GetRedisClient(config.URL, config.Password)
		client := NewClient(redisClient)
		return client, CloseRedis(client), err
	}
	redisClient, err := GetRedisClusterClient(config.URL, config.Password)
	client := NewClusterClient(redisClient)
	return client, CloseRedis(client), err
}
