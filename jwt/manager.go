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

// SetExtraResolver 设置刷新时加载用户信息的回调。
// 支持在 Manager 创建后延迟设置，用于解决循环依赖问题。
//
// 典型用法：
//
//	// 1. 先创建 jwtMgr（不带 resolver）
//	jwtMgr, _ := jwt.NewManager(jwt.WithSecret(cfg.Secret))
//
//	// 2. 注册到容器
//	container.Register(jwtMgr)
//
//	// 3. 创建依赖 jwtMgr 的服务
//	userService := NewUserService(jwtMgr, userRepo)
//
//	// 4. 延迟注入 resolver
//	jwtMgr.SetExtraResolver(userService)
func (m *Manager) SetExtraResolver(resolver ExtraResolver) {
	m.opts.Resolver = resolver
}

// NewManager 创建 JWT 管理器。
//
// 至少需要配置 AccessSecret 或 RefreshSecret 中的一个：
//   - 仅配置 AccessSecret：可生成/解析访问令牌
//   - 仅配置 RefreshSecret：可生成/解析刷新令牌
//   - 同时配置两者：完整功能
//
// 可通过 WithSecret 同时设置两者，或分别使用 WithAccessSecret 和 WithRefreshSecret。
func NewManager(opts ...Option) (*Manager, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.AccessSecret == "" && o.RefreshSecret == "" {
		return nil, ErrInvalidSecret
	}

	return &Manager{opts: o}, nil
}

// GenerateAccessToken 生成访问令牌。
// userID: 用户唯一标识
// username: 用户名
// extra: 自定义扩展字段（如角色、权限等）
//
// 若未配置 AccessSecret，返回 ErrAccessSecretNotConfigured。
func (m *Manager) GenerateAccessToken(userID, username string, extra map[string]any) (string, error) {
	if m.opts.AccessSecret == "" {
		return "", ErrAccessSecretNotConfigured
	}
	return generateToken(
		m.opts.AccessSecret,
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
//
// 若未配置 RefreshSecret，返回 ErrRefreshSecretNotConfigured。
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	if m.opts.RefreshSecret == "" {
		return "", ErrRefreshSecretNotConfigured
	}
	return generateToken(
		m.opts.RefreshSecret,
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

// ParseToken 尝试解析令牌并返回 Claims。
// 根据已配置的密钥尝试解析，优先使用 AccessSecret。
// 推荐使用 ParseAccessToken 或 ParseRefreshToken 以明确令牌类型。
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	var lastErr error

	// 尝试用 AccessSecret 解析
	if m.opts.AccessSecret != "" {
		claims, err := parseToken(tokenString, m.opts.AccessSecret)
		if err == nil {
			return claims, nil
		}
		lastErr = err
	}

	// 若 RefreshSecret 存在且与 AccessSecret 不同，尝试解析
	if m.opts.RefreshSecret != "" && m.opts.RefreshSecret != m.opts.AccessSecret {
		claims, err := parseToken(tokenString, m.opts.RefreshSecret)
		if err == nil {
			return claims, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ErrInvalidToken
}

// ParseAccessToken 解析访问令牌，若令牌类型不是 access 则返回错误。
//
// 若未配置 AccessSecret，返回 ErrAccessSecretNotConfigured。
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	if m.opts.AccessSecret == "" {
		return nil, ErrAccessSecretNotConfigured
	}

	claims, err := parseToken(tokenString, m.opts.AccessSecret)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ParseRefreshToken 解析刷新令牌，若令牌类型不是 refresh 则返回错误。
//
// 若未配置 RefreshSecret，返回 ErrRefreshSecretNotConfigured。
func (m *Manager) ParseRefreshToken(tokenString string) (*Claims, error) {
	if m.opts.RefreshSecret == "" {
		return nil, ErrRefreshSecretNotConfigured
	}

	claims, err := parseToken(tokenString, m.opts.RefreshSecret)
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
