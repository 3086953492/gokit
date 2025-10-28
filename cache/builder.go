package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
)

// Builder 泛型缓存构建器，提供链式调用的简洁 API
type Builder[T any] struct {
	key       string
	ttl       time.Duration
	localTTL  time.Duration
	skipLocal bool
}

// New 创建一个新的缓存构建器
// 使用示例：
//
//	builder := cache.New[models.User]()
func New[T any]() *Builder[T] {
	return &Builder[T]{
		ttl:       5 * time.Minute, // 默认 Redis TTL 5 分钟
		localTTL:  time.Minute,     // 默认本地缓存 TTL 1 分钟
		skipLocal: false,
	}
}

// Key 设置缓存键
func (b *Builder[T]) Key(key string) *Builder[T] {
	b.key = key
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
	c := GetGlobalCache()
	if c == nil {
		return nil, false, ErrCacheNotInitialized
	}

	var value T
	err := c.Get(ctx, b.key, &value)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &value, true, nil
}

// Set 设置缓存值
// ttl 参数可选，如果不传则使用 Builder 配置的 TTL
func (b *Builder[T]) Set(ctx context.Context, value *T, ttl ...time.Duration) error {
	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}

	cacheTTL := b.ttl
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	}

	item := &cache.Item{
		Ctx:            ctx,
		Key:            b.key,
		Value:          value,
		TTL:            cacheTTL,
		SkipLocalCache: b.skipLocal,
	}

	return c.Set(item)
}

// GetOrSet 获取缓存值，如果不存在则执行计算函数并缓存结果
// 这是最常用的方法，用于实现 cache-aside 模式
func (b *Builder[T]) GetOrSet(ctx context.Context, fn func() (*T, error)) (*T, error) {
	c := GetGlobalCache()
	if c == nil {
		return nil, ErrCacheNotInitialized
	}

	var value T
	err := c.Once(&cache.Item{
		Ctx:            ctx,
		Key:            b.key,
		Value:          &value,
		TTL:            b.ttl,
		SkipLocalCache: b.skipLocal,
		Do: func(*cache.Item) (interface{}, error) {
			return fn()
		},
	})

	if err != nil {
		return nil, err
	}

	return &value, nil
}

// Delete 删除缓存
func (b *Builder[T]) Delete(ctx context.Context) error {
	c := GetGlobalCache()
	if c == nil {
		return ErrCacheNotInitialized
	}
	return c.Delete(ctx, b.key)
}

// Exists 检查缓存键是否存在
func (b *Builder[T]) Exists(ctx context.Context) (bool, error) {
	c := GetGlobalCache()
	if c == nil {
		return false, ErrCacheNotInitialized
	}

	exists := c.Exists(ctx, b.key)
	return exists, nil
}
