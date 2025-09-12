package middleware

import (
	"github.com/3086953492/YaBase/middleware/oauth"
	"github.com/3086953492/YaBase/middleware/scopes"
	"github.com/3086953492/YaBase/middleware/security"
	"github.com/3086953492/YaBase/config"
	"github.com/3086953492/YaBase/config/types"

	"github.com/gin-gonic/gin"
)

// 中间件管理器 - 先做个简单版本
type Manager struct {
	config *types.MiddlewareConfig
}

// 创建管理器
func NewManager() *Manager {

	config := config.GetGlobalConfig().Middleware

	return &Manager{
		config: &config,
	}
}

// 加载所有全局中间件
func (m *Manager) LoadGlobal(engine *gin.Engine) {
	// CORS中间件
	engine.Use(m.CORS())

	// 添加Recovery中间件防止panic
	engine.Use(gin.Recovery())

	// 这里以后可以加更多全局中间件
	// engine.Use(m.Logger())
}

// CORS中间件
func (m *Manager) CORS() gin.HandlerFunc {
	return security.NewCORSMiddleware(security.CORSConfig{
		AllowOrigins: m.config.CORS.AllowOrigins,
		AllowMethods: m.config.CORS.AllowMethods,
		AllowHeaders: m.config.CORS.AllowHeaders,
	})
}

// 获取 OAuth 中间件
func (m *Manager) OAuth(requiredScopes ...string) gin.HandlerFunc {
	return oauth.OAuthTokenMiddleware(requiredScopes...)
}

func (m *Manager) RequiredScopes(scopeStrings ...string) gin.HandlerFunc {
	return scopes.RequiredScopes(scopeStrings...)
}
