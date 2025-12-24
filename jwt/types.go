package jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 定义令牌类型。
type TokenType string

const (
	// AccessToken 表示访问令牌，用于接口认证。
	AccessToken TokenType = "access"

	// RefreshToken 表示刷新令牌，用于换取新的访问令牌。
	RefreshToken TokenType = "refresh"
)

// Claims 定义 JWT 声明结构。
// 使用标准的 Subject (sub) 字段存储用户标识，并支持自定义扩展字段。
type Claims struct {
	// TokenType 令牌类型：access 或 refresh。
	TokenType TokenType `json:"token_type"`

	// Extra 自定义扩展字段，可存放角色、权限等信息。
	Extra map[string]any `json:"extra,omitempty"`

	jwt.RegisteredClaims
}

// ExtraResolver 定义刷新令牌时重新加载用户扩展信息的回调接口。
// 实现方应根据 subject 从数据库或缓存加载最新的扩展信息。
type ExtraResolver interface {
	// ResolveExtra 根据 subject 加载扩展信息。
	// 返回的 extra 将被写入新生成的 access token。
	ResolveExtra(ctx context.Context, subject string) (extra map[string]any, err error)
}

// ExtraResolverFunc 是 ExtraResolver 的函数适配器，
// 方便直接使用匿名函数作为 resolver。
type ExtraResolverFunc func(ctx context.Context, subject string) (map[string]any, error)

// ResolveExtra 实现 ExtraResolver 接口。
func (f ExtraResolverFunc) ResolveExtra(ctx context.Context, subject string) (map[string]any, error) {
	return f(ctx, subject)
}
