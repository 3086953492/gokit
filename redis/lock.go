package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DistributedLock 分布式锁实现
type DistributedLock struct {
	manager *Manager
	key     string
	value   string
	expire  time.Duration
}

// NewDistributedLock 创建一个新的分布式锁
func (m *Manager) NewDistributedLock(key string, expire time.Duration) *DistributedLock {
	return &DistributedLock{
		manager: m,
		key:     key,
		value:   uuid.New().String(),
		expire:  expire,
	}
}

// Acquire 尝试获取锁
// 如果获取成功返回 nil，否则返回 ErrLockAcquireFailed 或其他错误
func (l *DistributedLock) Acquire(ctx context.Context) error {
	ok, err := l.manager.SetNX(ctx, l.key, l.value, l.expire)
	if err != nil {
		return err
	}
	if !ok {
		return ErrLockAcquireFailed
	}
	return nil
}

// Release 释放锁
// 使用 Lua 脚本保证原子性，只有持有锁的进程才能释放
func (l *DistributedLock) Release(ctx context.Context) error {
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	_, err := l.manager.Eval(ctx, luaScript, []string{l.key}, l.value)
	return err
}
