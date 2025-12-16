package cache

import (
	"context"
	"time"
)

// RedisBackend 定义缓存所需的 Redis 操作接口。
// redis.Manager 自动满足此接口，无需显式实现。
type RedisBackend interface {
	// GetBytes 获取指定 key 的值
	// 如果 key 不存在，返回 (nil, nil)
	GetBytes(ctx context.Context, key string) ([]byte, error)

	// SetBytes 设置指定 key 的值
	SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Del 删除指定的 key，返回删除的数量
	Del(ctx context.Context, keys ...string) (int64, error)

	// ScanKeys 扫描匹配 pattern 的 key
	ScanKeys(ctx context.Context, pattern string, count int64) ([]string, error)

	// Exists 检查 key 是否存在
	Exists(ctx context.Context, key string) (bool, error)
}

// localCacheEntry 本地缓存条目
type localCacheEntry struct {
	value     []byte
	expireAt  time.Time
}

// isExpired 检查条目是否过期
func (e *localCacheEntry) isExpired() bool {
	return time.Now().After(e.expireAt)
}

