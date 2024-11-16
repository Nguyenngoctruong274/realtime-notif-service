package cachefx

import (
	"go.uber.org/fx"
	"yes4all/ads-noti-api/pkg/cache"
	"yes4all/ads-noti-api/pkg/xredis"
)

var Module = fx.Provide(provideRedisCache)

func provideRedisCache(redisCmd xredis.Cmdable) cache.RedisCache {
	return cache.New(redisCmd)
}
