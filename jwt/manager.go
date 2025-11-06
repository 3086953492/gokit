package jwt

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/3086953492/gokit/config/types"
)

// Manager JWT管理器
type Manager struct {
	config types.JWTConfig
}

// 全局JWT管理器
var (
	globalManager *Manager
	initOnce      sync.Once
	initErr       error
)

// InitJWT 初始化JWT管理器
// 此函数使用 sync.Once 保证只初始化一次
func InitJWT(cfg types.JWTConfig) error {
	initOnce.Do(func() {
		if cfg.Secret == "" {
			initErr = ErrInvalidToken
			return
		}
		globalManager = &Manager{
			config: cfg,
		}
	})
	return initErr
}

// GetGlobalJWT 获取全局JWT管理器实例
// 如果未初始化，返回 nil
func GetGlobalJWT() *Manager {
	return globalManager
}

// IsJWTInitialized 检查JWT管理器是否已初始化
func IsJWTInitialized() bool {
	return globalManager != nil
}

// GenerateToken 生成访问令牌
// userID: 用户ID
// username: 用户名
// extra: 额外的自定义字段
func (m *Manager) GenerateToken(userID, username string, extra map[string]interface{}) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: AccessToken,
		Extra:     extra,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.Expire)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// GenerateRefreshToken 生成刷新令牌
// userID: 用户ID
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshExpire)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// ParseToken 解析令牌并返回Claims
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.config.Secret), nil
	})

	if err != nil {
		// 检查是否是过期错误
		if jwt.ErrTokenExpired.Error() == err.Error() {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateToken 验证令牌是否有效
func (m *Manager) ValidateToken(tokenString string) error {
	_, err := m.ParseToken(tokenString)
	return err
}

// RefreshAccessToken 使用刷新令牌生成新的访问令牌
func (m *Manager) RefreshAccessToken(refreshToken string) (string, error) {
	// 解析刷新令牌
	claims, err := m.ParseToken(refreshToken)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	// 验证令牌类型
	if claims.TokenType != RefreshToken {
		return "", ErrInvalidTokenType
	}

	// 生成新的访问令牌
	return m.GenerateToken(claims.UserID, claims.Username, claims.Extra)
}

// GenerateToken 生成访问令牌（全局函数）
func GenerateToken(userID, username string, extra map[string]interface{}) (string, error) {
	m := GetGlobalJWT()
	if m == nil {
		return "", ErrJWTNotInitialized
	}
	return m.GenerateToken(userID, username, extra)
}

// GenerateRefreshToken 生成刷新令牌（全局函数）
func GenerateRefreshToken(userID string) (string, error) {
	m := GetGlobalJWT()
	if m == nil {
		return "", ErrJWTNotInitialized
	}
	return m.GenerateRefreshToken(userID)
}

// ParseToken 解析令牌（全局函数）
func ParseToken(tokenString string) (*Claims, error) {
	m := GetGlobalJWT()
	if m == nil {
		return nil, ErrJWTNotInitialized
	}
	return m.ParseToken(tokenString)
}

// ValidateToken 验证令牌（全局函数）
func ValidateToken(tokenString string) error {
	m := GetGlobalJWT()
	if m == nil {
		return ErrJWTNotInitialized
	}
	return m.ValidateToken(tokenString)
}

// RefreshAccessToken 刷新访问令牌（全局函数）
func RefreshAccessToken(refreshToken string) (string, error) {
	m := GetGlobalJWT()
	if m == nil {
		return "", ErrJWTNotInitialized
	}
	return m.RefreshAccessToken(refreshToken)
}

