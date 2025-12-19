package ginx

import "net/http"

// RedirectOptions 重定向配置
type RedirectOptions struct {
	Status int // HTTP 状态码，默认 302
}

func defaultRedirectOptions() *RedirectOptions {
	return &RedirectOptions{
		Status: http.StatusFound, // 302
	}
}

// RedirectOption 重定向选项函数
type RedirectOption func(*RedirectOptions)

// WithRedirectStatus 设置重定向状态码
func WithRedirectStatus(status int) RedirectOption {
	return func(o *RedirectOptions) {
		if status >= 300 && status < 400 {
			o.Status = status
		}
	}
}
