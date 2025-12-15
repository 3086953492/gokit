package storage

import "errors"

// 常见的存储错误。
// 调用方可以使用 errors.Is 判断具体错误类型。
var (
	// ErrNotFound 表示对象不存在。
	ErrNotFound = errors.New("storage: object not found")

	// ErrAlreadyExists 表示对象已存在（用于不允许覆盖的场景）。
	ErrAlreadyExists = errors.New("storage: object already exists")

	// ErrInvalidKey 表示对象 Key 不合法。
	ErrInvalidKey = errors.New("storage: invalid object key")

	// ErrBackendUnavailable 表示存储后端不可用。
	ErrBackendUnavailable = errors.New("storage: backend unavailable")

	// ErrPermissionDenied 表示无权限访问。
	ErrPermissionDenied = errors.New("storage: permission denied")

	// ErrInvalidConfig 表示配置无效。
	ErrInvalidConfig = errors.New("storage: invalid configuration")

	// ErrInvalidURL 表示提供的 URL 格式不合法或无法解析。
	ErrInvalidURL = errors.New("storage: invalid url")

	// ErrDomainNotAllowed 表示 URL 的域名不在当前 Store 允许范围内。
	ErrDomainNotAllowed = errors.New("storage: domain not allowed")

	// ErrURLDeleteUnsupported 表示当前 Store 不支持按 URL 删除。
	ErrURLDeleteUnsupported = errors.New("storage: delete by url not supported")
)
