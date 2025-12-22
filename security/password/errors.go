package password

import "errors"

// 密码哈希相关错误变量
var (
	// ErrEmptyPassword 表示密码为空
	ErrEmptyPassword = errors.New("密码不能为空")

	// ErrPasswordTooLong 表示密码超过最大长度限制
	ErrPasswordTooLong = errors.New("密码长度超过限制")

	// ErrInvalidCost 表示加密强度参数无效
	ErrInvalidCost = errors.New("加密强度参数无效")

	// ErrMismatch 表示密码与哈希不匹配
	ErrMismatch = errors.New("密码不匹配")

	// ErrHashInvalid 表示哈希值格式无效
	ErrHashInvalid = errors.New("哈希值格式无效")

	// ErrHashFailed 表示哈希生成失败
	ErrHashFailed = errors.New("密码哈希生成失败")
)

