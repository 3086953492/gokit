package initialize

import (
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func InitCache(redis *redis.Client) (*cache.Cache, error) {
	cacheInstance := cache.New(&cache.Options{
		Redis:      redis,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var isInitCache bool

	err := cacheInstance.Once(&cache.Item{
		Key:   "init",
		Value: &isInitCache,
		Do: func(*cache.Item) (any, error) {
			isInitCache = true
			return isInitCache, nil
		},
	})
	if err != nil || !isInitCache {
		return nil, err
	}
	return cacheInstance, nil
}
