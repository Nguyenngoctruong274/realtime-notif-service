package xredis

import (
	"sort"
	"sync"
	"yes4all/ads-noti-api/pkg/logger"

	"github.com/go-redis/redis"
)

type Cmdable interface {
	redis.UniversalClient
	AllKeys(matchKey string) *StringSliceCmd
}

type clusterClient struct {
	*redis.ClusterClient
}

func NewClusterClient(client *redis.ClusterClient) Cmdable {
	return &clusterClient{
		client,
	}
}

func (c *clusterClient) AllKeys(matchKey string) (result *StringSliceCmd) {
	keys := make([]string, 0, maxKeyCount)
	mutex := sync.Mutex{}
	err := c.ForEachMaster(func(client *redis.Client) error {
		sliceResult := c.scanKeys(client, matchKey)
		if sliceResult.err != nil {
			return sliceResult.err
		}
		mutex.Lock()
		defer mutex.Unlock()
		keys = append(keys, sliceResult.val...)

		return nil
	})

	if err != nil {
		return &StringSliceCmd{
			err: err,
		}
	}
	sort.Strings(keys)
	return &StringSliceCmd{
		val: keys,
	}
}

func (c *clusterClient) scanKeys(client *redis.Client, matchKey string) (result StringSliceCmd) {
	var count int64 = maxKeyCount
	var err error
	keys := make([]string, 0, count)
	var scanResult []string
	var cursor uint64
	for {
		scanResult, cursor, err = client.Scan(cursor, matchKey, count).Result()
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

func (c *clusterClient) Keys(pattern string) *redis.StringSliceCmd {
	logger.NewLogger().Warn("Warning: Keys method is return keys in one node only, please use AllKeys instead")
	return c.ClusterClient.Keys(pattern)
}

func (c *clusterClient) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	logger.NewLogger().Warnf("Warning: Scan method is return keys in one node only, please use AllKeys instead")
	return c.ClusterClient.Scan(cursor, match, count)
}
