package ginx

import "net/http"

// RedirectOptions 重定向配置
type RedirectOptions struct {
	Status int               // HTTP 状态码，默认 302
	Query  map[string]string // 附加 query 参数（同名 key 覆盖原有）
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

// WithRedirectQuery 设置附加的 query 参数，同名 key 会覆盖原有值
func WithRedirectQuery(q map[string]string) RedirectOption {
	return func(o *RedirectOptions) {
		o.Query = q
	}
}
