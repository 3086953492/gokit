package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/3086953492/gokit/config/types"
	"github.com/redis/go-redis/v9"
)

// 全局Redis管理器
var (
	globalRedis      *redis.Client
	redisMutex       sync.RWMutex
	redisInitialized bool
)

// GetGlobalRedis 获取全局Redis客户端
func GetGlobalRedis() *redis.Client {
	redisMutex.RLock()
	defer redisMutex.RUnlock()
	return globalRedis
}

// IsRedisInitialized 检查Redis是否已初始化
func IsRedisInitialized() bool {
	redisMutex.RLock()
	defer redisMutex.RUnlock()
	return redisInitialized
}

// InitRedisWithConfig 使用配置初始化全局Redis客户端
func InitRedisWithConfig(cfg types.RedisConfig) error {
	// 创建Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis连接失败: %v", err)
	}

	// 设置全局Redis客户端
	redisMutex.Lock()
	globalRedis = redisClient
	redisInitialized = true
	redisMutex.Unlock()

	fmt.Printf("Redis初始化成功: %s:%d\n", cfg.Host, cfg.Port)
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	redisMutex.Lock()
	defer redisMutex.Unlock()

	if globalRedis != nil {
		err := globalRedis.Close()
		globalRedis = nil
		redisInitialized = false
		return err
	}
	return nil
}
