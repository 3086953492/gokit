package storage

import (
	"context"
	"io"
	"time"
)

// Store 对象存储的统一后端接口。
// 所有具体存储实现（本地文件系统、阿里云 OSS、AWS S3 等）都应实现此接口。
type Store interface {
	// Upload 上传对象到存储后端。
	// key 为对象唯一标识，r 为对象内容流，opts 为可选写入参数。
	// 返回上传后的对象元信息。
	Upload(ctx context.Context, key string, r io.Reader, opts *WriteOptions) (*ObjectMeta, error)

	// Download 从存储后端下载对象。
	// 返回对象内容流（调用方负责 Close）和元信息。
	Download(ctx context.Context, key string, opts *ReadOptions) (io.ReadCloser, *ObjectMeta, error)

	// Delete 删除存储后端中的对象。
	Delete(ctx context.Context, key string, opts *DeleteOptions) error

	// List 列举指定前缀下的对象。
	// 返回的切片非 nil（即使为空）。
	List(ctx context.Context, prefix string, opts *ListOptions) (*ListResult, error)

	// Exists 检查对象是否存在。
	Exists(ctx context.Context, key string) (bool, error)

	// Head 获取对象元信息，不下载内容。
	Head(ctx context.Context, key string) (*ObjectMeta, error)
}

// ObjectMeta 对象元信息。
type ObjectMeta struct {
	Key          string            // 对象唯一标识
	Size         int64             // 对象大小（字节）
	ContentType  string            // MIME 类型
	ETag         string            // 对象 ETag（通常为内容的 MD5）
	LastModified time.Time         // 最后修改时间
	UserMeta     map[string]string // 自定义元数据
	URL          string            // 可访问直链（公开，不带签名）
}

// ListResult 列举结果。
type ListResult struct {
	Objects        []*ObjectMeta // 对象列表
	NextMarker     string        // 用于分页的下一页标记，空表示没有更多
	IsTruncated    bool          // 是否还有更多结果
	CommonPrefixes []string      // 当使用 Delimiter 时返回的公共前缀
}

// WriteOptions 上传时的可选参数。
type WriteOptions struct {
	ContentType   string            // MIME 类型
	CacheControl  string            // Cache-Control 头
	ContentLength int64             // 内容长度（某些后端需要预先知道）
	UserMeta      map[string]string // 自定义元数据
}

// ReadOptions 下载时的可选参数。
type ReadOptions struct {
	Range string // 范围读取，格式如 "bytes=0-1023"
}

// DeleteOptions 删除时的可选参数。
type DeleteOptions struct {
	// 暂无特殊选项，预留扩展
}

// ListOptions 列举时的可选参数。
type ListOptions struct {
	MaxKeys   int    // 单次返回的最大对象数，默认 1000
	Marker    string // 分页标记，从上一次 ListResult.NextMarker 获取
	Delimiter string // 目录分隔符，通常为 "/"
}
