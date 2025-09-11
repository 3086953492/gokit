package oauth

import (
	"slices"

	"github.com/3086953492/YaBase/config"

	"github.com/golang-jwt/jwt/v5"
)

// ParseOAuthToken 解析 OAuth 令牌
func ParseOAuthToken(tokenString string) (*OAuthClaims, error) {
	cfg := config.GetGlobalConfig()

	token, err := jwt.ParseWithClaims(tokenString, &OAuthClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*OAuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// ValidateScopes 验证权限范围
func ValidateScopes(requestedScopes, allowedScopes []string) bool {
	for _, requested := range requestedScopes {
		found := slices.Contains(allowedScopes, requested)
		if !found {
			return false
		}
	}
	return true
}
