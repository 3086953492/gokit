package oauth

import (
	"fmt"
	"time"

	"github.com/3086953492/YaBase/config"

	"github.com/golang-jwt/jwt/v5"
)

func cfg() *config.Config {
	return config.GetGlobalConfig()
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userID uint, clientID string, scopes []string) (string, error) {
	claims := OAuthClaims{
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		TokenType: "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg().OAuth.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg().OAuth.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg().JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, clientID string, scopes []string) (string, error) {
	claims := OAuthClaims{
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		TokenType: "refresh_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg().OAuth.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg().OAuth.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg().JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
