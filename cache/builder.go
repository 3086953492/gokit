package cache

import (
	"context"
	"time"
)

// Builder 泛型缓存构建器，提供链式调用的简洁 API。
// Builder 必须通过 Manager.Builder() 创建。
type Builder[T any] struct {
	manager   *Manager
	key       string
	ttl       time.Duration
	localTTL  time.Duration
	skipLocal bool
}

// Builder 创建一个新的缓存构建器
// 使用示例：
//
//	builder := manager.Builder[models.User]()
func (m *Manager) Builder() *Builder[any] {
	return &Builder[any]{
		manager:   m,
		ttl:       m.opts.DefaultTTL,
		localTTL:  m.opts.LocalCacheTTL,
		skipLocal: false,
	}
}

// NewBuilder 创建一个指定类型的缓存构建器
// 使用示例：
//
//	builder := cache.NewBuilder[models.User](manager)
func NewBuilder[T any](m *Manager) *Builder[T] {
	return &Builder[T]{
		manager:   m,
		ttl:       m.opts.DefaultTTL,
		localTTL:  m.opts.LocalCacheTTL,
		skipLocal: false,
	}
}

// Key 设置缓存键
func (b *Builder[T]) Key(key string) *Builder[T] {
	b.key = key
	return b
}

// KeyWithConds 使用前缀和条件 map 构造稳定的缓存键
// 内部会对条件进行标准化处理，保证相同条件生成相同的 key
// 示例：KeyWithConds("oauth_client", map[string]any{"id": 1})
func (b *Builder[T]) KeyWithConds(prefix string, conds map[string]any) *Builder[T] {
	b.key = BuildKeyFromConds(prefix, conds)
	return b
}

// KeyWithParts 使用前缀和多个部分构造缓存键
// 内部会对每个部分进行标准化处理
// 示例：KeyWithParts("user", userId, "profile")
func (b *Builder[T]) KeyWithParts(prefix string, parts ...any) *Builder[T] {
	b.key = BuildKey(prefix, parts...)
	return b
}

// TTL 设置 Redis 缓存过期时间
func (b *Builder[T]) TTL(ttl time.Duration) *Builder[T] {
	b.ttl = ttl
	return b
}

// LocalTTL 设置本地缓存过期时间
func (b *Builder[T]) LocalTTL(ttl time.Duration) *Builder[T] {
	b.localTTL = ttl
	return b
}

// SkipLocalCache 跳过本地缓存，只使用 Redis
func (b *Builder[T]) SkipLocalCache() *Builder[T] {
	b.skipLocal = true
	return b
}

// Get 从缓存中获取值
// 返回值、是否命中缓存、错误
func (b *Builder[T]) Get(ctx context.Context) (*T, bool, error) {
	var value T
	err := b.manager.Get(ctx, b.key, &value)
	if err != nil {
		if err == ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &value, true, nil
}

// Set 设置缓存值
// ttl 参数可选，如果不传则使用 Builder 配置的 TTL
func (b *Builder[T]) Set(ctx context.Context, value *T, ttl ...time.Duration) error {
	cacheTTL := b.ttl
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	}
	return b.manager.Set(ctx, b.key, value, cacheTTL)
}

// GetOrSet 获取缓存值，如果不存在则执行计算函数并缓存结果
// 这是最常用的方法，用于实现 cache-aside 模式
func (b *Builder[T]) GetOrSet(ctx context.Context, fn func() (*T, error)) (*T, error) {
	var value T
	err := b.manager.GetOrSet(ctx, b.key, &value, func() (any, error) {
		return fn()
	}, b.ttl)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// Delete 删除缓存
func (b *Builder[T]) Delete(ctx context.Context) error {
	return b.manager.Delete(ctx, b.key)
}

// Exists 检查缓存键是否存在
func (b *Builder[T]) Exists(ctx context.Context) (bool, error) {
	return b.manager.Exists(ctx, b.key)
}
