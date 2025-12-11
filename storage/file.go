package storage

import "io"

// FileObject 上传时的文件对象，调用方只需构造此结构
type FileObject struct {
	Reader      io.Reader // 文件内容流
	Size        int64     // 文件大小（字节）
	Filename    string    // 原始文件名（可选，用于推断扩展名）
	ContentType string    // MIME 类型（可选）
}

// UploadResult 上传结果
type UploadResult struct {
	Provider string // 实际使用的存储服务商（如 "aliyun_oss"、"s3" 等）
	Bucket   string // 实际使用的 Bucket 名称
	Key      string // 对象 Key
	URL      string // 可访问的完整 URL
}
