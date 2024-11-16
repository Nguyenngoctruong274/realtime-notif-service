package cache

import (
	"time"
	"yes4all/ads-noti-api/pkg/utils/json"
)

type Cache interface {
	Get(key string, data interface{}) error
	Set(key string, data interface{}) error
	SetWithExpiration(key string, data interface{}, expiration time.Duration) error
	Delete(key string) error
	AllKeys(pattern string) ([]string, error)
}

type RedisCache interface {
	Cache
	HGet(key, field string, data interface{}) error
	HGetAll(key string) (map[string]string, error)
	HSet(key, field string, data interface{}) error
	HDel(key string, fields ...string) error
	Incr(key string) (result int64, err error)
	Decr(key string) (result int64, err error)
	SetNX(key string, data interface{}) (result bool, err error)
}

// KeyFunc defines a transformer for cache keys
type KeyFunc func(s string) string

// DefaultKeyFunc is the default implementation of cache keys
// All it does is to return the key sent in by client code
func DefaultKeyFunc(s string) string {
	return s
}

type Serializer interface {
	Serialize(data interface{}) (string, error)
	Deserialize(data string, target interface{}) error
}

type JSONSerializer struct{}

func (s JSONSerializer) Serialize(data interface{}) (string, error) {
	return json.Serialize(data)
}

func (s JSONSerializer) Deserialize(data string, target interface{}) error {
	return json.Deserialize(data, target)
}
