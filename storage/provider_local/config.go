// Package providerlocal implements a local filesystem backend for storage.
package providerlocal

import "os"

// Config 本地文件系统存储配置。
type Config struct {
	Root     string      // 本地存储根目录，必填
	BaseURL  string      // 可选：公开访问前缀，用于生成 ObjectMeta.URL
	DirPerm  os.FileMode // 可选：自动创建目录时使用的权限，默认 0o755
	FilePerm os.FileMode // 可选：写入文件时使用的权限，默认 0o644
}
