// Package cookie 提供用于管理 access/refresh token Cookie 的辅助函数。
package cookie

import "net/http"

// 默认 Cookie 名称。
const (
	DefaultAccessName  = "access_token"
	DefaultRefreshName = "refresh_token"
	DefaultPath        = "/"
)

// Options Cookie 配置。
type Options struct {
	// AccessName 访问令牌 Cookie 名称，默认 "access_token"。
	AccessName string

	// RefreshName 刷新令牌 Cookie 名称，默认 "refresh_token"。
	RefreshName string

	// Domain Cookie 作用域，默认空（浏览器自动使用当前域名）。
	Domain string

	// Path Cookie 有效路径，默认 "/"。
	Path string

	// Secure 是否仅通过 HTTPS 传输，默认 false（生产环境建议开启）。
	Secure bool

	// HttpOnly 禁止 JS 访问，默认 true。
	HttpOnly bool

	// SameSite Cookie 同站策略，默认 Lax。
	SameSite http.SameSite

	// AccessMaxAge 访问令牌有效时长（秒），默认 900（15 分钟）。
	AccessMaxAge int

	// RefreshMaxAge 刷新令牌有效时长（秒），默认 604800（7 天）。
	RefreshMaxAge int
}

// DefaultOptions 返回带有合理默认值的 Options。
func DefaultOptions() *Options {
	return &Options{
		AccessName:    DefaultAccessName,
		RefreshName:   DefaultRefreshName,
		Path:          DefaultPath,
		Secure:        false,
		HttpOnly:      true,
		SameSite:      http.SameSiteLaxMode,
		AccessMaxAge:  900,
		RefreshMaxAge: 604800,
	}
}

// Option 配置函数类型。
type Option func(*Options)

// WithAccessName 设置访问令牌 Cookie 名称。
func WithAccessName(name string) Option {
	return func(o *Options) {
		if name != "" {
			o.AccessName = name
		}
	}
}

// WithRefreshName 设置刷新令牌 Cookie 名称。
func WithRefreshName(name string) Option {
	return func(o *Options) {
		if name != "" {
			o.RefreshName = name
		}
	}
}

// WithDomain 设置 Cookie 作用域。
func WithDomain(domain string) Option {
	return func(o *Options) {
		o.Domain = domain
	}
}

// WithPath 设置 Cookie 有效路径。
func WithPath(path string) Option {
	return func(o *Options) {
		if path != "" {
			o.Path = path
		}
	}
}

// WithSecure 设置是否仅通过 HTTPS 传输。
func WithSecure(secure bool) Option {
	return func(o *Options) {
		o.Secure = secure
	}
}

// WithHttpOnly 设置是否禁止 JS 访问。
func WithHttpOnly(httpOnly bool) Option {
	return func(o *Options) {
		o.HttpOnly = httpOnly
	}
}

// WithSameSite 设置 Cookie 同站策略。
func WithSameSite(sameSite http.SameSite) Option {
	return func(o *Options) {
		o.SameSite = sameSite
	}
}

// WithAccessMaxAge 设置访问令牌有效时长（秒）。
func WithAccessMaxAge(maxAge int) Option {
	return func(o *Options) {
		o.AccessMaxAge = maxAge
	}
}

// WithRefreshMaxAge 设置刷新令牌有效时长（秒）。
func WithRefreshMaxAge(maxAge int) Option {
	return func(o *Options) {
		o.RefreshMaxAge = maxAge
	}
}

