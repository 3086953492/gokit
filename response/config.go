package response

import (
	"sync"

	"github.com/3086953492/gokit/errors"
)

// Config 响应包的配置
type Config struct {
	// ShowErrorDetail 是否显示错误详细信息（包括 Cause 和 Fields）
	ShowErrorDetail bool

	// ErrorStatusMapper 错误类型到 HTTP 状态码的映射函数
	// 如果为 nil，则使用默认映射
	ErrorStatusMapper func(errType string) int

	// OAuthErrorCodeMapper 错误类型到 OAuth 2.0 错误码的映射函数
	// 如果为 nil，则使用默认映射
	OAuthErrorCodeMapper func(errType string) string

	// DefaultSuccessMessage 默认成功消息
	DefaultSuccessMessage string

	// DefaultPaginatedMessage 默认分页成功消息
	DefaultPaginatedMessage string

	// FallbackErrorMessage 兜底错误消息（非 AppError 时使用）
	FallbackErrorMessage string

	// FallbackErrorCode 兜底错误代码（非 AppError 时使用）
	FallbackErrorCode string
}

// Option 配置选项函数
type Option func(*Config)

var (
	globalConfig *Config
	configMutex  sync.RWMutex
	initOnce     sync.Once
)

// Init 初始化响应包配置
// 此函数应在应用启动时调用一次，多次调用只有第一次生效
func Init(opts ...Option) {
	initOnce.Do(func() {
		cfg := defaultConfig()
		for _, opt := range opts {
			opt(cfg)
		}
		configMutex.Lock()
		globalConfig = cfg
		configMutex.Unlock()
	})
}

// getConfig 获取全局配置，如果未初始化则返回默认配置
func getConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if globalConfig == nil {
		return defaultConfig()
	}
	return globalConfig
}

// defaultConfig 返回默认配置
func defaultConfig() *Config {
	return &Config{
		ShowErrorDetail:         false,
		ErrorStatusMapper:       defaultErrorStatusMapper,
		OAuthErrorCodeMapper:    defaultOAuthErrorCodeMapper,
		DefaultSuccessMessage:   "操作成功",
		DefaultPaginatedMessage: "获取成功",
		FallbackErrorMessage:    "系统内部错误",
		FallbackErrorCode:       "SYSTEM_ERROR",
	}
}

// defaultErrorStatusMapper 默认的错误类型到 HTTP 状态码映射
func defaultErrorStatusMapper(errType string) int {
	switch errType {
	case errors.TypeNotFound:
		return 404
	case errors.TypeInvalidInput:
		return 400
	case errors.TypeUnauthorized:
		return 401
	case errors.TypeForbidden:
		return 403
	case errors.TypeDuplicate:
		return 409
	case errors.TypeValidation:
		return 422
	default:
		return 500
	}
}

// defaultOAuthErrorCodeMapper 默认的错误类型到 OAuth 2.0 错误码映射
func defaultOAuthErrorCodeMapper(errType string) string {
	switch errType {
	case errors.TypeInvalidInput:
		return "invalid_request"
	case errors.TypeUnauthorized:
		return "unauthorized_client"
	case errors.TypeForbidden:
		return "access_denied"
	default:
		return "server_error"
	}
}

// WithShowErrorDetail 设置是否显示错误详细信息
func WithShowErrorDetail(show bool) Option {
	return func(c *Config) {
		c.ShowErrorDetail = show
	}
}

// WithErrorStatusMapper 自定义错误类型到 HTTP 状态码的映射
func WithErrorStatusMapper(mapper func(errType string) int) Option {
	return func(c *Config) {
		if mapper != nil {
			c.ErrorStatusMapper = mapper
		}
	}
}

// WithOAuthErrorCodeMapper 自定义错误类型到 OAuth 2.0 错误码的映射
func WithOAuthErrorCodeMapper(mapper func(errType string) string) Option {
	return func(c *Config) {
		if mapper != nil {
			c.OAuthErrorCodeMapper = mapper
		}
	}
}

// WithDefaultSuccessMessage 设置默认成功消息
func WithDefaultSuccessMessage(msg string) Option {
	return func(c *Config) {
		c.DefaultSuccessMessage = msg
	}
}

// WithDefaultPaginatedMessage 设置默认分页成功消息
func WithDefaultPaginatedMessage(msg string) Option {
	return func(c *Config) {
		c.DefaultPaginatedMessage = msg
	}
}

// WithFallbackErrorMessage 设置兜底错误消息
func WithFallbackErrorMessage(msg string) Option {
	return func(c *Config) {
		c.FallbackErrorMessage = msg
	}
}

// WithFallbackErrorCode 设置兜底错误代码
func WithFallbackErrorCode(code string) Option {
	return func(c *Config) {
		c.FallbackErrorCode = code
	}
}
