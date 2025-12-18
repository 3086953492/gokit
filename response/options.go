package response

import "net/http"

// Options Manager 的配置选项。
type Options struct {
	// DefaultSuccessMessage 成功响应的默认消息。
	DefaultSuccessMessage string
	// DefaultPaginatedMessage 分页成功响应的默认消息。
	DefaultPaginatedMessage string

	// FallbackStatus 无法从 Error 获取状态码时使用的默认状态码。
	FallbackStatus int
	// FallbackCode 无法从 error 解析出 Error 时使用的默认错误码。
	FallbackCode string
	// FallbackTitle 兜底 title（对应 FallbackCode 时使用）。
	FallbackTitle string

	// ProblemTypePrefix Problem type 前缀；默认 "urn:problem-type:"。
	// 最终 type 为 ProblemTypePrefix + code。
	ProblemTypePrefix string

	// OAuth2ErrorMapper 将错误映射为 OAuth2 参数 (error, error_description)。
	// 返回 ok=false 表示不注入 OAuth2 错误参数。
	OAuth2ErrorMapper func(err error) (code string, desc string, ok bool)
}

// Option 配置函数类型。
type Option func(*Options)

// defaultOptions 返回默认配置。
func defaultOptions() *Options {
	return &Options{
		DefaultSuccessMessage:   "操作成功",
		DefaultPaginatedMessage: "获取成功",
		FallbackStatus:          http.StatusInternalServerError,
		FallbackCode:            "INTERNAL_ERROR",
		FallbackTitle:           "Internal Server Error",
		ProblemTypePrefix:       "urn:problem-type:",
		OAuth2ErrorMapper:       nil,
	}
}

// WithDefaultSuccessMessage 设置成功响应的默认消息。
func WithDefaultSuccessMessage(msg string) Option {
	return func(o *Options) {
		o.DefaultSuccessMessage = msg
	}
}

// WithDefaultPaginatedMessage 设置分页成功响应的默认消息。
func WithDefaultPaginatedMessage(msg string) Option {
	return func(o *Options) {
		o.DefaultPaginatedMessage = msg
	}
}

// WithFallbackStatus 设置无法解析状态码时的默认状态码。
func WithFallbackStatus(status int) Option {
	return func(o *Options) {
		o.FallbackStatus = status
	}
}

// WithFallbackCode 设置无法解析 Error 时使用的默认错误码。
func WithFallbackCode(code string) Option {
	return func(o *Options) {
		o.FallbackCode = code
	}
}

// WithFallbackTitle 设置兜底 title。
func WithFallbackTitle(title string) Option {
	return func(o *Options) {
		o.FallbackTitle = title
	}
}

// WithProblemTypePrefix 设置 problem type 的前缀。
func WithProblemTypePrefix(prefix string) Option {
	return func(o *Options) {
		o.ProblemTypePrefix = prefix
	}
}

// WithOAuth2ErrorMapper 设置 OAuth2 错误映射函数。
func WithOAuth2ErrorMapper(mapper func(err error) (code string, desc string, ok bool)) Option {
	return func(o *Options) {
		o.OAuth2ErrorMapper = mapper
	}
}

