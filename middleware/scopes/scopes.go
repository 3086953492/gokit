package scopes

import (
	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/response"
	"github.com/gin-gonic/gin"
	"slices"
)

// Scopes权限中间件
func RequiredScopes(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userScopes, exists := c.Get("scopes")
		if !exists {
			response.Error(c, apperrors.ErrInsufficientScope)
			c.Abort()
			return
		}

		scopes, ok := userScopes.([]string)
		if !ok {
			response.Error(c, apperrors.ErrInsufficientScope)
			c.Abort()
			return
		}

		// 检查是否有所需的权限范围
		for _, required := range requiredScopes {
			found := slices.Contains(scopes, required)
			if !found {
				response.Error(c, apperrors.ErrInsufficientScope)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
