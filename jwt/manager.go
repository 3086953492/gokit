package jwt

import (
	"context"
	"fmt"
)

// Manager 提供 JWT 令牌的生成、解析与刷新功能。
// Manager 是线程安全的，可在多个 goroutine 中共享使用。
type Manager struct {
	opts *Options
}

// NewManager 创建 JWT 管理器。
// 必须通过 WithSecret 指定签名密钥，否则返回 ErrInvalidSecret。
func NewManager(opts ...Option) (*Manager, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.Secret == "" {
		return nil, ErrInvalidSecret
	}

	return &Manager{opts: o}, nil
}

// GenerateAccessToken 生成访问令牌。
// userID: 用户唯一标识
// username: 用户名
// extra: 自定义扩展字段（如角色、权限等）
func (m *Manager) GenerateAccessToken(userID, username string, extra map[string]any) (string, error) {
	return generateToken(
		m.opts.Secret,
		m.opts.Issuer,
		m.opts.AccessTTL,
		AccessToken,
		userID,
		username,
		extra,
	)
}

// GenerateRefreshToken 生成刷新令牌。
// 刷新令牌仅包含 userID，不携带敏感信息。
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	return generateToken(
		m.opts.Secret,
		m.opts.Issuer,
		m.opts.RefreshTTL,
		RefreshToken,
		userID,
		"",  // 不写入 username
		nil, // 不写入 extra
	)
}

// GenerateTokenPair 同时生成访问令牌和刷新令牌。
// 返回 (accessToken, refreshToken, error)。
func (m *Manager) GenerateTokenPair(userID, username string, extra map[string]any) (string, string, error) {
	accessToken, err := m.GenerateAccessToken(userID, username, extra)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析令牌并返回 Claims。
// 支持解析 access 和 refresh 两种类型的令牌。
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, m.opts.Secret)
}

// ParseAccessToken 解析访问令牌，若令牌类型不是 access 则返回错误。
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ParseRefreshToken 解析刷新令牌，若令牌类型不是 refresh 则返回错误。
func (m *Manager) ParseRefreshToken(tokenString string) (*Claims, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ValidateToken 验证令牌是否有效。
func (m *Manager) ValidateToken(tokenString string) error {
	_, err := m.ParseToken(tokenString)
	return err
}

// RefreshAccessToken 使用刷新令牌生成新的访问令牌。
// 此方法需要配置 ExtraResolver，用于根据 userID 加载最新的用户信息。
// 若未配置 resolver，返回 ErrResolverNotConfigured。
func (m *Manager) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// 检查 resolver 是否已配置
	if m.opts.Resolver == nil {
		return "", ErrResolverNotConfigured
	}

	// 解析刷新令牌
	claims, err := m.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 通过 resolver 加载最新的用户信息
	username, extra, err := m.opts.Resolver.ResolveExtra(ctx, claims.UserID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrResolveFailed, err)
	}

	// 生成新的访问令牌
	return m.GenerateAccessToken(claims.UserID, username, extra)
}
