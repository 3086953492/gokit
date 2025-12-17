package jwt

import (
	"context"
	"time"
)

// 默认配置值。
const (
	// DefaultAccessTTL 默认访问令牌有效期：15 分钟。
	DefaultAccessTTL = 15 * time.Minute

	// DefaultRefreshTTL 默认刷新令牌有效期：7 天。
	DefaultRefreshTTL = 7 * 24 * time.Hour

	// DefaultIssuer 默认签发者。
	DefaultIssuer = "gokit"
)

// Options 包含 Manager 的配置参数。
type Options struct {
	// Secret 签名密钥，必填。
	Secret string

	// Issuer 令牌签发者。
	Issuer string

	// AccessTTL 访问令牌有效期。
	AccessTTL time.Duration

	// RefreshTTL 刷新令牌有效期。
	RefreshTTL time.Duration

	// Resolver 刷新时用于加载用户信息的回调，可选。
	// 若未配置，调用 RefreshAccessToken 将返回 ErrResolverNotConfigured。
	Resolver ExtraResolver
}

// Option 是配置 Manager 的函数类型。
type Option func(*Options)

// defaultOptions 返回带有默认值的 Options。
func defaultOptions() *Options {
	return &Options{
		Issuer:     DefaultIssuer,
		AccessTTL:  DefaultAccessTTL,
		RefreshTTL: DefaultRefreshTTL,
	}
}

// WithSecret 设置签名密钥。
func WithSecret(secret string) Option {
	return func(o *Options) {
		o.Secret = secret
	}
}

// WithIssuer 设置令牌签发者。
func WithIssuer(issuer string) Option {
	return func(o *Options) {
		o.Issuer = issuer
	}
}

// WithAccessTTL 设置访问令牌有效期。
func WithAccessTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.AccessTTL = ttl
	}
}

// WithRefreshTTL 设置刷新令牌有效期。
func WithRefreshTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.RefreshTTL = ttl
	}
}

// WithExtraResolver 设置刷新时加载用户信息的回调。
func WithExtraResolver(resolver ExtraResolver) Option {
	return func(o *Options) {
		o.Resolver = resolver
	}
}

// WithExtraResolverFunc 使用函数作为 ExtraResolver。
func WithExtraResolverFunc(fn func(ctx context.Context, userID string) (string, map[string]any, error)) Option {
	return func(o *Options) {
		o.Resolver = ExtraResolverFunc(fn)
	}
}

