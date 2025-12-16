package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// Manager 缓存管理器，提供两级缓存（本地缓存 + Redis）。
// Manager 是线程安全的，可以在多个 goroutine 中共享使用。
type Manager struct {
	redis RedisBackend
	opts  *Options
	local *localCache
	sf    singleflight.Group

	mu     sync.RWMutex
	closed bool
}

// NewManager 创建一个新的缓存 Manager。
// redis 参数是 RedisBackend 接口的实现（如 redis.Manager）。
func NewManager(redis RedisBackend, opts ...Option) (*Manager, error) {
	if redis == nil {
		return nil, ErrNilRedisBackend
	}

	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	var local *localCache
	if options.LocalCacheEnabled {
		local = newLocalCache(options.LocalCacheMaxSize)
	}

	return &Manager{
		redis: redis,
		opts:  options,
		local: local,
	}, nil
}

// Close 关闭 Manager
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}
	m.closed = true

	if m.local != nil {
		m.local.clear()
	}
	return nil
}

// Get 从缓存获取值并反序列化到 dest
// 如果缓存未命中，返回 ErrCacheMiss
func (m *Manager) Get(ctx context.Context, key string, dest any) error {
	if err := m.checkClosed(); err != nil {
		return err
	}

	// 先查本地缓存
	if m.local != nil {
		if data := m.local.get(key); data != nil {
			return json.Unmarshal(data, dest)
		}
	}

	// 查 Redis
	data, err := m.redis.GetBytes(ctx, key)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrCacheMiss
	}

	// 写入本地缓存
	if m.local != nil {
		m.local.set(key, data, m.opts.LocalCacheTTL)
	}

	return json.Unmarshal(data, dest)
}

// Set 序列化值并写入缓存
func (m *Manager) Set(ctx context.Context, key string, value any, ttl ...time.Duration) error {
	if err := m.checkClosed(); err != nil {
		return err
	}

	cacheTTL := m.opts.DefaultTTL
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}

	// 写入 Redis
	if err := m.redis.SetBytes(ctx, key, data, cacheTTL); err != nil {
		return err
	}

	// 写入本地缓存
	if m.local != nil {
		m.local.set(key, data, m.opts.LocalCacheTTL)
	}

	return nil
}

// GetOrSet 获取缓存值，如果不存在则执行 fn 并缓存结果
// 使用 singleflight 防止缓存击穿
func (m *Manager) GetOrSet(ctx context.Context, key string, dest any, fn func() (any, error), ttl ...time.Duration) error {
	if err := m.checkClosed(); err != nil {
		return err
	}

	// 先尝试获取
	err := m.Get(ctx, key, dest)
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		return err
	}

	// 使用 singleflight 防止并发请求
	result, err, _ := m.sf.Do(key, func() (any, error) {
		// 双重检查
		if err := m.Get(ctx, key, dest); err == nil {
			return nil, nil
		}

		// 执行回调
		value, err := fn()
		if err != nil {
			return nil, err
		}

		// 写入缓存
		if err := m.Set(ctx, key, value, ttl...); err != nil {
			return nil, err
		}

		return value, nil
	})

	if err != nil {
		return err
	}

	// 如果 result 不为 nil，说明是本次调用执行了 fn
	if result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("cache marshal: %w", err)
		}
		return json.Unmarshal(data, dest)
	}

	return nil
}

// Delete 删除指定 key 的缓存
func (m *Manager) Delete(ctx context.Context, key string) error {
	if err := m.checkClosed(); err != nil {
		return err
	}

	// 删除本地缓存
	if m.local != nil {
		m.local.delete(key)
	}

	// 删除 Redis
	_, err := m.redis.Del(ctx, key)
	return err
}

// DeleteByPrefix 删除指定前缀的所有缓存
func (m *Manager) DeleteByPrefix(ctx context.Context, prefix string) error {
	if err := m.checkClosed(); err != nil {
		return err
	}

	// 删除本地缓存
	if m.local != nil {
		m.local.deleteByPrefix(prefix)
	}

	// 扫描并删除 Redis 中的 key
	keys, err := m.redis.ScanKeys(ctx, prefix+"*", m.opts.ScanCount)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		_, err = m.redis.Del(ctx, keys...)
		return err
	}

	return nil
}

// DeleteByContains 删除包含指定子串的所有缓存
func (m *Manager) DeleteByContains(ctx context.Context, substring string) error {
	if err := m.checkClosed(); err != nil {
		return err
	}

	// 扫描并删除 Redis 中的 key
	keys, err := m.redis.ScanKeys(ctx, "*"+substring+"*", m.opts.ScanCount)
	if err != nil {
		return err
	}

	// 删除本地缓存中匹配的 key
	if m.local != nil {
		for _, key := range keys {
			m.local.delete(key)
		}
	}

	if len(keys) > 0 {
		_, err = m.redis.Del(ctx, keys...)
		return err
	}

	return nil
}

// DeleteByPrefixes 批量删除多个前缀的所有缓存
func (m *Manager) DeleteByPrefixes(ctx context.Context, prefixes []string) error {
	if len(prefixes) == 0 {
		return nil
	}

	var errs []error
	for _, prefix := range prefixes {
		if err := m.DeleteByPrefix(ctx, prefix); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Exists 检查缓存键是否存在
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	if err := m.checkClosed(); err != nil {
		return false, err
	}

	// 先查本地缓存
	if m.local != nil {
		if data := m.local.get(key); data != nil {
			return true, nil
		}
	}

	// 查 Redis
	return m.redis.Exists(ctx, key)
}

// GetKeysByPrefix 获取指定前缀的所有缓存键
func (m *Manager) GetKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	if err := m.checkClosed(); err != nil {
		return nil, err
	}

	return m.redis.ScanKeys(ctx, prefix+"*", m.opts.ScanCount)
}

// GetKeysByContains 获取包含指定子串的所有缓存键
func (m *Manager) GetKeysByContains(ctx context.Context, substring string) ([]string, error) {
	if err := m.checkClosed(); err != nil {
		return nil, err
	}

	return m.redis.ScanKeys(ctx, "*"+substring+"*", m.opts.ScanCount)
}

// DeleteByConds 根据前缀和条件删除缓存
// 使用标准化的 key 生成规则，确保能准确删除对应缓存
func (m *Manager) DeleteByConds(ctx context.Context, prefix string, conds map[string]any) error {
	key := BuildKeyFromConds(prefix, conds)
	return m.Delete(ctx, key)
}

// DeleteByCondsPrefix 删除指定前缀+条件前缀的所有缓存
// 例如：DeleteByCondsPrefix(ctx, "oauth_client", map[string]any{"id": 1})
// 会删除所有以 "oauth_client|id=1" 开头的缓存
func (m *Manager) DeleteByCondsPrefix(ctx context.Context, prefix string, conds map[string]any) error {
	keyPrefix := BuildKeyFromConds(prefix, conds)
	return m.DeleteByPrefix(ctx, keyPrefix)
}

// DeleteByContainsList 根据前缀和多组条件批量删除缓存
// 每组条件会通过 BuildKeyFromConds 标准化为稳定的子串，然后使用 DeleteByContains 逻辑删除
// 示例：DeleteByContainsList(ctx, "oauth_client", []map[string]any{{"id": 1}, {"id": 2}})
func (m *Manager) DeleteByContainsList(ctx context.Context, prefix string, condsList []map[string]any) error {
	if len(condsList) == 0 {
		return nil
	}

	var errs []error
	for _, conds := range condsList {
		substring := BuildKeyFromConds(prefix, conds)
		if err := m.DeleteByContains(ctx, substring); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// checkClosed 检查 Manager 是否已关闭
func (m *Manager) checkClosed() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.closed {
		return ErrManagerClosed
	}
	return nil
}
