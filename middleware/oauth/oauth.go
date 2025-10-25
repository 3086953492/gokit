package oauth

import (
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"

	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/response"
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

		token := tokenParts[1]

		claims, err := casdoorsdk.ParseJwtToken(token)
		if err != nil {
			panic(err)
		}

		c.Set("user", claims.User)
        
		c.Next()
	}
}
