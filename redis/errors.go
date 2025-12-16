package redis

import "errors"

var (
	// ErrNotConnected 表示 Manager 尚未连接到 Redis
	ErrNotConnected = errors.New("redis: not connected")

	// ErrAlreadyConnected 表示 Manager 已连接
	ErrAlreadyConnected = errors.New("redis: already connected")

	// ErrClosed 表示 Manager 已关闭
	ErrClosed = errors.New("redis: manager closed")

	// ErrLockAcquireFailed 表示获取分布式锁失败
	ErrLockAcquireFailed = errors.New("redis: lock acquire failed")
)

