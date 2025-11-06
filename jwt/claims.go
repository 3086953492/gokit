package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// TokenType 令牌类型
type TokenType string

const (
	// AccessToken 访问令牌
	AccessToken TokenType = "access"
	// RefreshToken 刷新令牌
	RefreshToken TokenType = "refresh"
)

// Claims JWT声明结构
type Claims struct {
	UserID    string         `json:"user_id"`
	Username  string         `json:"username,omitempty"`
	TokenType TokenType      `json:"token_type"`
	Extra     map[string]any `json:"extra,omitempty"`
	jwt.RegisteredClaims
}
