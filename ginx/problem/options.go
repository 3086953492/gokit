package problem

// Options Fail 函数配置
type Options struct {
	Instance   string         // 覆盖 Problem.Instance
	Extensions map[string]any // 扩展字段，序列化时展平到顶层
}

func defaultOptions() *Options {
	return &Options{}
}

// Option 失败响应选项函数
type Option func(*Options)

// WithInstance 设置 Problem.Instance（通常为请求 URI）
func WithInstance(instance string) Option {
	return func(o *Options) {
		o.Instance = instance
	}
}

// WithExtensions 设置扩展字段，序列化时展平到 Problem 顶层
func WithExtensions(ext map[string]any) Option {
	return func(o *Options) {
		o.Extensions = ext
	}
}

