package response

// Options 成功响应配置
type Options struct {
	Message string // 成功消息，默认 "ok"
}

func defaultOptions() *Options {
	return &Options{
		Message: "ok",
	}
}

// Option 成功响应选项函数
type Option func(*Options)

// WithMessage 设置自定义成功消息
func WithMessage(msg string) Option {
	return func(o *Options) {
		if msg != "" {
			o.Message = msg
		}
	}
}

