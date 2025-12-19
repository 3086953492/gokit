// Package redirect 提供 HTTP 重定向辅助函数。
package redirect

import "net/http"

// Options 重定向配置
type Options struct {
	Status int               // HTTP 状态码，默认 302
	Query  map[string]string // 附加 query 参数（同名 key 覆盖原有）
}

func defaultOptions() *Options {
	return &Options{
		Status: http.StatusFound, // 302
	}
}

// Option 重定向选项函数
type Option func(*Options)

// WithStatus 设置重定向状态码
func WithStatus(status int) Option {
	return func(o *Options) {
		if status >= 300 && status < 400 {
			o.Status = status
		}
	}
}

// WithQuery 设置附加的 query 参数，同名 key 会覆盖原有值
func WithQuery(q map[string]string) Option {
	return func(o *Options) {
		o.Query = q
	}
}

