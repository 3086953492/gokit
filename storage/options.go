package storage

// Options Manager 的配置选项。
type Options struct {
	store Store
}

// Option Manager 配置函数类型。
type Option func(*Options)

// defaultOptions 返回默认配置。
func defaultOptions() *Options {
	return &Options{}
}

// WithStore 设置存储后端实现。
// 这是必需选项，Manager 初始化时必须指定。
func WithStore(store Store) Option {
	return func(o *Options) {
		o.store = store
	}
}

// WriteOption 上传时的可选配置函数类型。
type WriteOption func(*WriteOptions)

// WithContentType 设置对象的 MIME 类型。
func WithContentType(ct string) WriteOption {
	return func(o *WriteOptions) {
		o.ContentType = ct
	}
}

// WithCacheControl 设置对象的 Cache-Control 头。
func WithCacheControl(cc string) WriteOption {
	return func(o *WriteOptions) {
		o.CacheControl = cc
	}
}

// WithContentLength 设置对象的内容长度。
// 某些后端（如阿里云 OSS）需要预先知道内容长度。
func WithContentLength(length int64) WriteOption {
	return func(o *WriteOptions) {
		o.ContentLength = length
	}
}

// WithUserMeta 设置自定义元数据。
func WithUserMeta(meta map[string]string) WriteOption {
	return func(o *WriteOptions) {
		o.UserMeta = meta
	}
}

// ReadOption 下载时的可选配置函数类型。
type ReadOption func(*ReadOptions)

// WithRange 设置范围读取。
// 格式如 "bytes=0-1023"。
func WithRange(r string) ReadOption {
	return func(o *ReadOptions) {
		o.Range = r
	}
}

// DeleteOption 删除时的可选配置函数类型。
type DeleteOption func(*DeleteOptions)

// ListOption 列举时的可选配置函数类型。
type ListOption func(*ListOptions)

// WithMaxKeys 设置单次返回的最大对象数。
func WithMaxKeys(n int) ListOption {
	return func(o *ListOptions) {
		o.MaxKeys = n
	}
}

// WithMarker 设置分页标记。
func WithMarker(marker string) ListOption {
	return func(o *ListOptions) {
		o.Marker = marker
	}
}

// WithDelimiter 设置目录分隔符。
// 通常为 "/"，用于模拟目录结构。
func WithDelimiter(delimiter string) ListOption {
	return func(o *ListOptions) {
		o.Delimiter = delimiter
	}
}
