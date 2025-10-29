package cache

import (
	"context"
	"fmt"

	"github.com/3086953492/gokit/redis"
	"github.com/go-redis/cache/v9"
	redislib "github.com/redis/go-redis/v9"
)

// getCacheKeysByPrefix 获取指定前缀的所有缓存键（私有函数）
func getCacheKeysByPrefix(ctx context.Context, prefix string, redisClient *redislib.Client) ([]string, error) {
	var keys []string
	pattern := prefix + "*"

	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("扫描缓存键失败: %w", err)
	}

	return keys, nil
}

// deleteCacheKeysByPrefix 删除指定前缀的所有缓存键（私有函数）
func deleteCacheKeysByPrefix(ctx context.Context, prefix string, redisClient *redislib.Client, cacheClient *cache.Cache) error {
	keys, err := getCacheKeysByPrefix(ctx, prefix, redisClient)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil // 没有匹配的键，直接返回
	}

	// 批量删除缓存
	for _, key := range keys {
		if err := cacheClient.Delete(ctx, key); err != nil {
			return fmt.Errorf("删除缓存键 %s 失败: %w", key, err)
		}
	}

	return nil
}

// GetKeysByPrefix 获取指定前缀的所有缓存键（公开函数）
func GetKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	redisClient := redis.GetGlobalRedis()
	if redisClient == nil {
		return nil, ErrRedisNotInitialized
	}

	return getCacheKeysByPrefix(ctx, prefix, redisClient)
}
