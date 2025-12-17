package jwt

import "errors"

// 哨兵错误定义，调用方可使用 errors.Is 判断具体错误类型。
var (
	// ErrInvalidSecret 表示提供的签名密钥为空或无效。
	ErrInvalidSecret = errors.New("jwt: invalid secret")

	// ErrInvalidToken 表示令牌格式不正确或签名验证失败。
	ErrInvalidToken = errors.New("jwt: invalid token")

	// ErrTokenExpired 表示令牌已过期。
	ErrTokenExpired = errors.New("jwt: token expired")

	// ErrInvalidTokenType 表示令牌类型与预期不符（如期望 refresh 却传入 access）。
	ErrInvalidTokenType = errors.New("jwt: invalid token type")

	// ErrResolverNotConfigured 表示未配置 ExtraResolver，无法执行刷新操作。
	ErrResolverNotConfigured = errors.New("jwt: extra resolver not configured")

	// ErrResolveFailed 表示 ExtraResolver 执行失败。
	ErrResolveFailed = errors.New("jwt: resolve extra failed")
)
