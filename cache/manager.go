package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"

	"github.com/3086953492/gokit/redis"
)

var (
	// ErrCacheNotInitialized 缓存未初始化错误
	ErrCacheNotInitialized = errors.New("缓存未初始化，请先调用 InitCache()")

	// ErrRedisNotInitialized Redis 未初始化错误
	ErrRedisNotInitialized = errors.New("redis 客户端未初始化")
)

// 全局缓存管理器
var (
	globalCache *cache.Cache
	initOnce    sync.Once
	initErr     error
)

// GetGlobalCache 获取全局缓存实例
// 如果缓存未初始化，返回 nil
func GetGlobalCache() *cache.Cache {
	return globalCache
}

// IsCacheInitialized 检查缓存是否已初始化
func IsCacheInitialized() bool {
	return globalCache != nil
}

// InitCache 使用 Redis 初始化全局缓存
// 注意：需要先初始化 Redis 模块，然后再调用此函数
// 此函数使用 sync.Once 保证只初始化一次
func InitCache() error {
	initOnce.Do(func() {
		// 获取 Redis 客户端
		redisClient := redis.GetGlobalRedis()
		if redisClient == nil {
			initErr = ErrRedisNotInitialized
			return
		}

		// 初始化缓存（本地缓存 + Redis 两级缓存）
		globalCache = cache.New(&cache.Options{
			Redis:      redisClient,
			LocalCache: cache.NewTinyLFU(1000, time.Minute), // 本地缓存 1000 个条目，默认 1 分钟过期
		})
	})

	return initErr
}

// Delete 删除指定键的缓存
func Delete(ctx context.Context, key string) error {
	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}
	return c.Delete(ctx, key)
}

// DeleteByPrefix 删除指定前缀的所有缓存
func DeleteByPrefix(ctx context.Context, prefix string) error {
	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}

	redisClient := redis.GetGlobalRedis()
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	return deleteCacheKeysByPrefix(ctx, prefix, redisClient, c)
}

// DeleteByContains 删除包含指定子串的所有缓存
func DeleteByContains(ctx context.Context, substring string) error {
	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}

	redisClient := redis.GetGlobalRedis()
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	return deleteCacheKeysByContains(ctx, substring, redisClient, c)
}

// DeleteByPrefixes 批量删除多个前缀的所有缓存
func DeleteByPrefixes(ctx context.Context, prefixes []string) error {
	if len(prefixes) == 0 {
		return nil
	}

	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}

	redisClient := redis.GetGlobalRedis()
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	var errs []error
	for _, prefix := range prefixes {
		if err := deleteCacheKeysByPrefix(ctx, prefix, redisClient, c); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// DeleteByContainsList 批量删除包含多个子串的所有缓存
func DeleteByContainsList(ctx context.Context, substrings []string) error {
	if len(substrings) == 0 {
		return nil
	}

	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}

	redisClient := redis.GetGlobalRedis()
	if redisClient == nil {
		return ErrRedisNotInitialized
	}

	var errs []error
	for _, substring := range substrings {
		if err := deleteCacheKeysByContains(ctx, substring, redisClient, c); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Exists 检查缓存键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	c := GetGlobalCache()
	if c == nil {
		return false, ErrCacheNotInitialized
	}
	exists := c.Exists(ctx, key)
	return exists, nil
}
