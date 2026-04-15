package cache

import (
	"context"
	"time"
)

// KeyedOption 配置 Keyed 实例。
type KeyedOption func(*keyedConfig)

type keyedConfig struct {
	ttl time.Duration
}

// WithKeyedTTL 设置 Keyed 的 Redis 缓存过期时间。
// 不设置时取 Manager 的 DefaultTTL。
func WithKeyedTTL(ttl time.Duration) KeyedOption {
	return func(c *keyedConfig) {
		c.ttl = ttl
	}
}

// Keyed 是预定义的类型化缓存访问器。
// 它将 Manager、键前缀、TTL 和值类型一次性绑定，
// 调用方只需提供可变的键组成部分（parts）。
//
// Keyed 是线程安全的，可在多个 goroutine 中共享使用。
//
// 使用示例：
//
//	listCache := cache.NewKeyed[dto.ProductList](mgr, "product:list", cache.WithKeyedTTL(10*time.Minute))
//	result, err := listCache.GetOrSet(ctx, fn, page, pageSize, lang)
type Keyed[T any] struct {
	mgr    *Manager
	prefix string
	ttl    time.Duration
}

// NewKeyed 创建预定义的类型化缓存访问器。
// prefix 作为缓存键的固定前缀，各方法的 parts 参数通过 BuildKey 追加在后面生成完整 key。
func NewKeyed[T any](mgr *Manager, prefix string, opts ...KeyedOption) *Keyed[T] {
	cfg := keyedConfig{ttl: mgr.opts.DefaultTTL}
	for _, o := range opts {
		o(&cfg)
	}
	return &Keyed[T]{mgr: mgr, prefix: prefix, ttl: cfg.ttl}
}

// Prefix 返回此访问器的键前缀，可用于构造 Group。
func (k *Keyed[T]) Prefix() string { return k.prefix }

// Get 按 parts 生成缓存键并获取值。
// 缓存未命中时返回 (nil, false, nil)。
func (k *Keyed[T]) Get(ctx context.Context, parts ...any) (*T, bool, error) {
	var value T
	err := k.mgr.Get(ctx, BuildKey(k.prefix, parts...), &value)
	if err != nil {
		if err == ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &value, true, nil
}

// Set 按 parts 生成缓存键并写入值，使用预配置的 TTL。
func (k *Keyed[T]) Set(ctx context.Context, value *T, parts ...any) error {
	return k.mgr.Set(ctx, BuildKey(k.prefix, parts...), value, k.ttl)
}

// GetOrSet 按 parts 生成缓存键，命中则返回，未命中则执行 fn 并缓存结果。
// 内部使用 singleflight 防止缓存击穿。
func (k *Keyed[T]) GetOrSet(ctx context.Context, fn func() (*T, error), parts ...any) (*T, error) {
	var value T
	key := BuildKey(k.prefix, parts...)
	err := k.mgr.GetOrSet(ctx, key, &value, func() (any, error) {
		return fn()
	}, k.ttl)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// Delete 按 parts 生成缓存键并删除对应缓存。
func (k *Keyed[T]) Delete(ctx context.Context, parts ...any) error {
	return k.mgr.Delete(ctx, BuildKey(k.prefix, parts...))
}

// InvalidateAll 删除此前缀下的所有缓存。
func (k *Keyed[T]) InvalidateAll(ctx context.Context) error {
	return k.mgr.DeleteByPrefix(ctx, k.prefix)
}

// Exists 按 parts 生成缓存键并检查是否存在。
func (k *Keyed[T]) Exists(ctx context.Context, parts ...any) (bool, error) {
	return k.mgr.Exists(ctx, BuildKey(k.prefix, parts...))
}
