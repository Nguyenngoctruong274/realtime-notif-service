package redisfx

import (
	"yes4all/ads-noti-api/pkg/infra"
	provider "yes4all/ads-noti-api/pkg/provider/redis"

	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Provide(
	provideRedisClient,
)

func provideRedisClient(redis infra.RedisClient) provider.IRedisProvider {
	return provider.NewRedisProvider(redis)
}
