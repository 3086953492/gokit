package cookie

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenCookies 提供 access/refresh token Cookie 的读写与清理能力。
type TokenCookies struct {
	opts *Options
}

// New 创建 TokenCookies 实例。
func New(opts ...Option) *TokenCookies {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &TokenCookies{opts: o}
}

// SetAccess 向响应中写入访问令牌 Cookie。
func (t *TokenCookies) SetAccess(c *gin.Context, token string) {
	t.setCookie(c, t.opts.AccessName, token, t.opts.AccessMaxAge)
}

// SetRefresh 向响应中写入刷新令牌 Cookie。
func (t *TokenCookies) SetRefresh(c *gin.Context, token string) {
	t.setCookie(c, t.opts.RefreshName, token, t.opts.RefreshMaxAge)
}

// GetAccess 从请求中读取访问令牌 Cookie。
func (t *TokenCookies) GetAccess(c *gin.Context) (string, error) {
	return c.Cookie(t.opts.AccessName)
}

// GetRefresh 从请求中读取刷新令牌 Cookie。
func (t *TokenCookies) GetRefresh(c *gin.Context) (string, error) {
	return c.Cookie(t.opts.RefreshName)
}

// Clear 清除 access 与 refresh 两个 Cookie（通过设置 MaxAge=-1）。
func (t *TokenCookies) Clear(c *gin.Context) {
	t.setCookie(c, t.opts.AccessName, "", -1)
	t.setCookie(c, t.opts.RefreshName, "", -1)
}

// setCookie 使用标准库 http.SetCookie 写入带 SameSite 属性的 Cookie。
func (t *TokenCookies) setCookie(c *gin.Context, name, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     t.opts.Path,
		Domain:   t.opts.Domain,
		MaxAge:   maxAge,
		Secure:   t.opts.Secure,
		HttpOnly: t.opts.HttpOnly,
		SameSite: t.opts.SameSite,
	})
}
