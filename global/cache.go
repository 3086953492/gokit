package global

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
)

// 全局缓存管理器
var (
	globalCache      *cache.Cache
	cacheMutex       sync.RWMutex
	cacheInitialized bool
)

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache() *cache.Cache {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return globalCache
}

// IsCacheInitialized 检查缓存是否已初始化
func IsCacheInitialized() bool {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return cacheInitialized
}

// InitCacheWithConfig 使用配置初始化全局缓存
// 注意：需要先初始化Redis模块，然后再调用此函数
func InitCacheWithConfig() error {

	// 获取Redis客户端
	redisClient := GetGlobalRedis()
	if redisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}

	// 初始化缓存
	cacheInstance := cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	// 设置全局缓存
	cacheMutex.Lock()
	globalCache = cacheInstance
	cacheInitialized = true
	cacheMutex.Unlock()

	return nil
}
