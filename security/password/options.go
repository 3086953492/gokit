package password

import "golang.org/x/crypto/bcrypt"

// Options 定义 Manager 的配置选项
type Options struct {
	// Cost bcrypt 加密强度，范围 4-31，默认 10
	Cost int

	// MaxPasswordBytes 密码最大字节数，默认 72（bcrypt 限制）
	MaxPasswordBytes int

	// AllowEmpty 是否允许空密码，默认 false
	AllowEmpty bool
}

// Option 是配置选项函数类型
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		Cost:             bcrypt.DefaultCost,
		MaxPasswordBytes: 72,
		AllowEmpty:       false,
	}
}

// WithCost 设置 bcrypt 加密强度
// cost 范围为 4-31，值越大越安全但计算时间越长
func WithCost(cost int) Option {
	return func(o *Options) {
		o.Cost = cost
	}
}

// WithMaxPasswordBytes 设置密码最大字节数
// bcrypt 最大支持 72 字节
func WithMaxPasswordBytes(max int) Option {
	return func(o *Options) {
		o.MaxPasswordBytes = max
	}
}

// WithAllowEmpty 设置是否允许空密码
func WithAllowEmpty(allow bool) Option {
	return func(o *Options) {
		o.AllowEmpty = allow
	}
}

