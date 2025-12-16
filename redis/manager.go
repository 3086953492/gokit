package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Manager 管理 Redis 连接，提供对 Redis 的操作接口。
// Manager 是线程安全的，可以在多个 goroutine 中共享使用。
type Manager struct {
	opts   *Options
	client *redis.Client

	mu       sync.RWMutex
	closed   bool
	connOnce sync.Once
	connErr  error
}

// NewManager 创建一个新的 Redis Manager。
// 调用 NewManager 后需要显式调用 Connect(ctx) 来建立连接。
func NewManager(opts ...Option) *Manager {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Manager{
		opts: options,
	}
}

// Connect 建立到 Redis 的连接。
// 此方法是幂等的，多次调用只会建立一次连接。
// 如果 Manager 已关闭，返回 ErrClosed。
func (m *Manager) Connect(ctx context.Context) error {
	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return ErrClosed
	}
	m.mu.RUnlock()

	m.connOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:         m.opts.Address,
			Password:     m.opts.Password,
			DB:           m.opts.DB,
			DialTimeout:  m.opts.DialTimeout,
			ReadTimeout:  m.opts.ReadTimeout,
			WriteTimeout: m.opts.WriteTimeout,
			PoolSize:     m.opts.PoolSize,
			MinIdleConns: m.opts.MinIdleConns,
		})

		if err := client.Ping(ctx).Err(); err != nil {
			m.connErr = fmt.Errorf("redis connect failed: %w", err)
			return
		}

		m.mu.Lock()
		m.client = client
		m.mu.Unlock()
	})

	return m.connErr
}

// Close 关闭 Redis 连接。
// 关闭后 Manager 不可再使用。
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}
	m.closed = true

	if m.client != nil {
		return m.client.Close()
	}
	return nil
}

// IsConnected 返回是否已连接到 Redis
func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client != nil && !m.closed
}

// getClient 获取底层 Redis 客户端，内部使用
func (m *Manager) getClient() (*redis.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, ErrClosed
	}
	if m.client == nil {
		return nil, ErrNotConnected
	}
	return m.client, nil
}

// GetBytes 获取指定 key 的值（字节切片形式）
func (m *Manager) GetBytes(ctx context.Context, key string) ([]byte, error) {
	client, err := m.getClient()
	if err != nil {
		return nil, err
	}

	result, err := client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get: %w", err)
	}
	return result, nil
}

// SetBytes 设置指定 key 的值（字节切片形式）
func (m *Manager) SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	client, err := m.getClient()
	if err != nil {
		return err
	}

	if err := client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

// Del 删除指定的 key，返回删除的 key 数量
func (m *Manager) Del(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	client, err := m.getClient()
	if err != nil {
		return 0, err
	}

	result, err := client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis del: %w", err)
	}
	return result, nil
}

// Exists 检查 key 是否存在
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	client, err := m.getClient()
	if err != nil {
		return false, err
	}

	result, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists: %w", err)
	}
	return result > 0, nil
}

// ScanKeys 扫描匹配 pattern 的 key
// pattern 支持通配符，如 "prefix*"、"*suffix"、"*contains*"
// count 是每次扫描的建议数量（实际返回可能更多或更少）
func (m *Manager) ScanKeys(ctx context.Context, pattern string, count int64) ([]string, error) {
	client, err := m.getClient()
	if err != nil {
		return nil, err
	}

	var keys []string
	iter := client.Scan(ctx, 0, pattern, count).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("redis scan: %w", err)
	}

	if keys == nil {
		keys = []string{}
	}
	return keys, nil
}

// SetNX 仅当 key 不存在时设置值，返回是否设置成功
func (m *Manager) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	client, err := m.getClient()
	if err != nil {
		return false, err
	}

	result, err := client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}
	return result, nil
}

// Eval 执行 Lua 脚本
func (m *Manager) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	client, err := m.getClient()
	if err != nil {
		return nil, err
	}

	result, err := client.Eval(ctx, script, keys, args...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis eval: %w", err)
	}
	return result, nil
}
