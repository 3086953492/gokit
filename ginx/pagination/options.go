package pagination

// Options 分页解析配置
type Options struct {
	DefaultPage     int
	DefaultPageSize int
	MaxPageSize     int
}

func defaultOptions() *Options {
	return &Options{
		DefaultPage:     1,
		DefaultPageSize: 10,
		MaxPageSize:     100,
	}
}

// Option 分页选项函数
type Option func(*Options)

// WithDefaultPage 设置默认页码
func WithDefaultPage(p int) Option {
	return func(o *Options) {
		if p > 0 {
			o.DefaultPage = p
		}
	}
}

// WithDefaultPageSize 设置默认每页条数
func WithDefaultPageSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.DefaultPageSize = size
		}
	}
}

// WithMaxPageSize 设置最大每页条数（防止客户端请求过大）
func WithMaxPageSize(max int) Option {
	return func(o *Options) {
		if max > 0 {
			o.MaxPageSize = max
		}
	}
}

