package cache

import "errors"

var (
	// ErrNilRedisBackend 表示 RedisBackend 为 nil
	ErrNilRedisBackend = errors.New("cache: redis backend is nil")

	// ErrInvalidKey 表示缓存 key 无效
	ErrInvalidKey = errors.New("cache: invalid key")

	// ErrCacheMiss 表示缓存未命中
	ErrCacheMiss = errors.New("cache: miss")

	// ErrManagerClosed 表示 Manager 已关闭
	ErrManagerClosed = errors.New("cache: manager closed")
)

