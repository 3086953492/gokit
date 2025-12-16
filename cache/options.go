package cache

import "time"

// Options 包含 Manager 的配置选项
type Options struct {
	// DefaultTTL 默认 Redis 缓存过期时间
	DefaultTTL time.Duration

	// LocalCacheEnabled 是否启用本地缓存
	LocalCacheEnabled bool

	// LocalCacheTTL 本地缓存过期时间
	LocalCacheTTL time.Duration

	// LocalCacheMaxSize 本地缓存最大条目数（0 表示不限制）
	LocalCacheMaxSize int

	// ScanCount 扫描 key 时每次迭代的建议数量
	ScanCount int64
}

// Option 是配置 Manager 的函数类型
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		DefaultTTL:        5 * time.Minute,
		LocalCacheEnabled: true,
		LocalCacheTTL:     time.Minute,
		LocalCacheMaxSize: 1000,
		ScanCount:         100,
	}
}

// WithDefaultTTL 设置默认 Redis 缓存过期时间
func WithDefaultTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.DefaultTTL = ttl
	}
}

// WithLocalCache 启用或禁用本地缓存
func WithLocalCache(enabled bool) Option {
	return func(o *Options) {
		o.LocalCacheEnabled = enabled
	}
}

// WithLocalCacheTTL 设置本地缓存过期时间
func WithLocalCacheTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.LocalCacheTTL = ttl
	}
}

// WithLocalCacheMaxSize 设置本地缓存最大条目数
func WithLocalCacheMaxSize(size int) Option {
	return func(o *Options) {
		o.LocalCacheMaxSize = size
	}
}

// WithScanCount 设置扫描 key 时每次迭代的建议数量
func WithScanCount(count int64) Option {
	return func(o *Options) {
		o.ScanCount = count
	}
}

