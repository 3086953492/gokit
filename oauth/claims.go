package oauth

import (
	"github.com/golang-jwt/jwt/v5"
)

// OAuthClaims OAuth JWT 声明结构
type OAuthClaims struct {
	UserID    uint     `json:"sub"`
	ClientID  string   `json:"client_id"`
	Scopes    []string `json:"scope"`
	TokenType string   `json:"token_type"` // access_token, refresh_token
	jwt.RegisteredClaims
}
