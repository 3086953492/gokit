package subject

import "errors"

// Subject 生成相关错误变量
var (
	// ErrEmptyUserID 表示用户 ID 为空
	ErrEmptyUserID = errors.New("用户 ID 不能为空")

	// ErrEmptySecret 表示密钥为空
	ErrEmptySecret = errors.New("密钥不能为空")

	// ErrInvalidLength 表示截断长度无效
	ErrInvalidLength = errors.New("截断长度无效")

	// ErrSecretTooShort 表示密钥长度过短（建议至少 32 字节）
	ErrSecretTooShort = errors.New("密钥长度过短，建议至少 32 字节")
)

