package errors

const (
	TypeSystemInfoNotFound     = "SYSTEM_INFO_NOT_FOUND"
	TypeSystemInfoListNotFound = "SYSTEM_INFO_LIST_NOT_FOUND"
)

var (
	ErrSystemInfoNotFound     = New(TypeSystemInfoNotFound, "系统配置不存在")
	ErrSystemInfoListNotFound = New(TypeSystemInfoListNotFound, "系统配置列表不存在")
)
