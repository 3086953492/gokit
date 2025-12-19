package ginx

// ---------------------------------------------------------------------------
// 分页选项
// ---------------------------------------------------------------------------

// PageOptions 分页解析配置
type PageOptions struct {
	DefaultPage     int
	DefaultPageSize int
	MaxPageSize     int
}

func defaultPageOptions() *PageOptions {
	return &PageOptions{
		DefaultPage:     1,
		DefaultPageSize: 10,
		MaxPageSize:     100,
	}
}

// PageOption 分页选项函数
type PageOption func(*PageOptions)

// WithDefaultPage 设置默认页码
func WithDefaultPage(p int) PageOption {
	return func(o *PageOptions) {
		if p > 0 {
			o.DefaultPage = p
		}
	}
}

// WithDefaultPageSize 设置默认每页条数
func WithDefaultPageSize(size int) PageOption {
	return func(o *PageOptions) {
		if size > 0 {
			o.DefaultPageSize = size
		}
	}
}

// WithMaxPageSize 设置最大每页条数（防止客户端请求过大）
func WithMaxPageSize(max int) PageOption {
	return func(o *PageOptions) {
		if max > 0 {
			o.MaxPageSize = max
		}
	}
}

// ---------------------------------------------------------------------------
// 失败响应选项
// ---------------------------------------------------------------------------

// FailOptions Fail 函数配置
type FailOptions struct {
	Instance   string         // 覆盖 Problem.Instance
	Extensions map[string]any // 扩展字段，序列化时展平到顶层
}

func defaultFailOptions() *FailOptions {
	return &FailOptions{}
}

// FailOption 失败响应选项函数
type FailOption func(*FailOptions)

// WithInstance 设置 Problem.Instance（通常为请求 URI）
func WithInstance(instance string) FailOption {
	return func(o *FailOptions) {
		o.Instance = instance
	}
}

// WithExtensions 设置扩展字段，序列化时展平到 Problem 顶层
func WithExtensions(ext map[string]any) FailOption {
	return func(o *FailOptions) {
		o.Extensions = ext
	}
}
