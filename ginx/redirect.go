package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Redirect 执行 HTTP 重定向并终止后续 handler 处理。
//
// 默认使用 302 Found，可通过 WithRedirectStatus 覆盖。
//
//	ginx.Redirect(c, "/login")
//	ginx.Redirect(c, "/dashboard", ginx.WithRedirectStatus(http.StatusMovedPermanently))
func Redirect(c *gin.Context, location string, opts ...RedirectOption) {
	o := defaultRedirectOptions()
	for _, fn := range opts {
		fn(o)
	}
	c.Redirect(o.Status, location)
	c.Abort()
}

// RedirectTo 执行指定状态码的 HTTP 重定向并终止后续 handler 处理。
//
//	ginx.RedirectTo(c, http.StatusMovedPermanently, "/new-path")
func RedirectTo(c *gin.Context, status int, location string) {
	c.Redirect(status, location)
	c.Abort()
}

// ---------------------------------------------------------------------------
// 便捷方法
// ---------------------------------------------------------------------------

// RedirectFound 执行 302 Found 重定向
func RedirectFound(c *gin.Context, location string) {
	RedirectTo(c, http.StatusFound, location)
}

// RedirectPermanent 执行 301 Moved Permanently 重定向
func RedirectPermanent(c *gin.Context, location string) {
	RedirectTo(c, http.StatusMovedPermanently, location)
}

// RedirectTemporary 执行 307 Temporary Redirect 重定向（保持请求方法和请求体）
func RedirectTemporary(c *gin.Context, location string) {
	RedirectTo(c, http.StatusTemporaryRedirect, location)
}

// RedirectSeeOther 执行 303 See Other 重定向（常用于 POST 后跳转到 GET）
func RedirectSeeOther(c *gin.Context, location string) {
	RedirectTo(c, http.StatusSeeOther, location)
}

