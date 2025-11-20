package cookie

import (
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
)

// 包级配置
type config struct {
	secure bool
}

var (
	cfg   = config{}
	cfgMu sync.RWMutex
)

// Init 初始化 cookie 包的配置。
// 应在应用启动时调用一次，用于设置是否仅通过 HTTPS 发送 Cookie。
func Init(secure bool) {
	cfgMu.Lock()
	cfg.secure = secure
	cfgMu.Unlock()
}

// internal helper：获取当前 secure 配置
func isSecure() bool {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg.secure
}

// SetAccessToken 设置访问令牌 Cookie。
// token: 令牌值
// maxAge: 有效时长
// domain, path: Cookie 作用域
func SetAccessToken(c *gin.Context, token string, maxAge int, domain, path string) {
	c.SetCookie(accessTokenCookieName, token, maxAge, path, domain, isSecure(), true)
}

// SetRefreshToken 设置刷新令牌 Cookie。
// token: 刷新令牌值
// maxAge: 有效时长
// domain, path: Cookie 作用域
func SetRefreshToken(c *gin.Context, token string, maxAge int, domain, path string) {
	c.SetCookie(refreshTokenCookieName, token, maxAge, path, domain, isSecure(), true)
}

// ClearTokens 清理访问令牌和刷新令牌 Cookie。
// 会通过设置过期时间的方式通知浏览器删除这两个 Cookie。
func ClearTokens(c *gin.Context, domain, path string) {
	clearCookie(c, accessTokenCookieName, domain, path)
	clearCookie(c, refreshTokenCookieName, domain, path)
}

// GetAccessToken 从 Cookie 中获取访问令牌。
func GetAccessToken(c *gin.Context) (string, error) {
	return c.Cookie(accessTokenCookieName)
}

// GetRefreshToken 从 Cookie 中获取刷新令牌。
func GetRefreshToken(c *gin.Context) (string, error) {
	return c.Cookie(refreshTokenCookieName)
}

// clearCookie 将指定名称的 Cookie 标记为过期以触发删除。
func clearCookie(c *gin.Context, name, domain, path string) {
	// maxAge 设为 -1 表示删除 Cookie。
	c.SetCookie(name, "", -1, path, domain, isSecure(), true)
}
