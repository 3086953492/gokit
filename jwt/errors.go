package jwt

import "errors"

var (
	// ErrJWTNotInitialized JWT管理器未初始化错误
	ErrJWTNotInitialized = errors.New("JWT管理器未初始化，请先调用 InitJWT()")

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = errors.New("无效的令牌")

	// ErrInvalidTokenType 无效的令牌类型
	ErrInvalidTokenType = errors.New("无效的令牌类型")
)
