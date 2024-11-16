package xredis

import "github.com/go-redis/redis"

type Client struct {
	*redis.Client
}

func NewClient(client *redis.Client) Cmdable {
	return &Client{
		Client: client,
	}
}

func (c *Client) AllKeys(pattern string) (result *StringSliceCmd) {
	result = &StringSliceCmd{}
	var count int64 = maxKeyCount
	var err error
	keys := make([]string, 0, count)
	var scanResult []string
	var cursor uint64
	for {
		scanResult, cursor, err = c.Client.Scan(cursor, pattern, count).Result()
		if err != nil {
			result.err = err
			return
		}
		keys = append(keys, scanResult...)
		if cursor == 0 {
			break
		}
	}

	result.val = keys
	return
}
