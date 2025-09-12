package oauth

import (
	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/oauth"
	"github.com/3086953492/YaBase/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// OAuth 访问令牌验证中间件
func OAuthTokenMiddleware(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 头中获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		accessToken := tokenParts[1]

		// 解析和验证令牌
		claims, err := oauth.ParseOAuthToken(accessToken)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// 检查令牌类型
		if claims.TokenType != "access_token" {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 检查权限范围
		if len(requiredScopes) > 0 && !oauth.ValidateScopes(requiredScopes, claims.Scopes) {
			response.Error(c, apperrors.ErrInsufficientScope)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("client_id", claims.ClientID)
		c.Set("scopes", claims.Scopes)
		c.Set("token_type", "oauth")

		c.Next()
	}
}
