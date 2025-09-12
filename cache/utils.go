package cache

import (
	"context"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func GetCacheKeysByPrefix(prefix string, redisClient *redis.Client) ([]string, error) {
	var keys []string
	iter := redisClient.Scan(context.Background(), 0, prefix+"*", 0).Iterator()
	for iter.Next(context.Background()) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func DeleteCacheKeysByPrefix(prefix string, redisClient *redis.Client, cacheClient *cache.Cache) error {
	keys, err := GetCacheKeysByPrefix(prefix, redisClient)
	if err != nil {
		return err
	}
	for _, key := range keys {
		err := cacheClient.Delete(context.Background(), key)
		if err != nil {
			return err
		}
	}
	return nil
}
