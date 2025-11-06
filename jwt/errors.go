package jwt

import "errors"

var (
	// ErrJWTNotInitialized JWT管理器未初始化错误
	ErrJWTNotInitialized = errors.New("JWT管理器未初始化，请先调用 InitJWT()")

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = errors.New("无效的令牌")

	// ErrExpiredToken 令牌已过期
	ErrExpiredToken = errors.New("令牌已过期")

	// ErrInvalidRefreshToken 无效的刷新令牌
	ErrInvalidRefreshToken = errors.New("无效的刷新令牌")

	// ErrInvalidTokenType 无效的令牌类型
	ErrInvalidTokenType = errors.New("无效的令牌类型")
)

