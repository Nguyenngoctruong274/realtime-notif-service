package cache

import (
	"log"
	"time"
)

type cacheOptions struct {
	expiration time.Duration
	keyFunc    KeyFunc
	serializer Serializer
}

// An Option configures a cache client.
type Option interface {
	apply(*cacheOptions)
}

type optionFunc func(*cacheOptions)

func (f optionFunc) apply(args *cacheOptions) {
	f(args)
}

func WithExpiration(e time.Duration) Option {
	return optionFunc(func(args *cacheOptions) {
		args.expiration = e
	})
}

// WithCacheKeyGenerator allows configuring the cache key generation process
func WithCacheKeyGenerator(fn KeyFunc) Option {
	if fn == nil {
		log.Fatal("KeyFunc cannot be nil")
	}
	return optionFunc(func(args *cacheOptions) {
		args.keyFunc = fn
	})
}

func WithSerializer(s Serializer) Option {
	return optionFunc(func(args *cacheOptions) {
		args.serializer = s
	})
}
