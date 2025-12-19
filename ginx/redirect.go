package ginx

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Redirect 执行 HTTP 重定向并终止后续 handler 处理。
//
// 默认使用 302 Found，可通过 WithRedirectStatus 覆盖。
// 可通过 WithRedirectQuery 附加 query 参数。
//
//	ginx.Redirect(c, "/login")
//	ginx.Redirect(c, "/dashboard", ginx.WithRedirectStatus(http.StatusMovedPermanently))
//	ginx.Redirect(c, "/callback", ginx.WithRedirectQuery(map[string]string{"token": "abc"}))
func Redirect(c *gin.Context, location string, opts ...RedirectOption) {
	o := defaultRedirectOptions()
	for _, fn := range opts {
		fn(o)
	}

	finalLocation := buildRedirectLocation(location, o.Query)
	c.Redirect(o.Status, finalLocation)
	c.Abort()
}

// RedirectTo 执行指定状态码的 HTTP 重定向并终止后续 handler 处理。
//
//	ginx.RedirectTo(c, http.StatusMovedPermanently, "/new-path")
func RedirectTo(c *gin.Context, status int, location string) {
	c.Redirect(status, location)
	c.Abort()
}

// RedirectToWithQuery 执行指定状态码的 HTTP 重定向并附加 query 参数。
//
//	ginx.RedirectToWithQuery(c, 302, "/callback", map[string]string{"code": "123"})
func RedirectToWithQuery(c *gin.Context, status int, location string, query map[string]string) {
	finalLocation := buildRedirectLocation(location, query)
	c.Redirect(status, finalLocation)
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

// ---------------------------------------------------------------------------
// 内部工具函数
// ---------------------------------------------------------------------------

// buildRedirectLocation 解析 location 并合并 query 参数（同名 key 覆盖）
func buildRedirectLocation(location string, query map[string]string) string {
	if len(query) == 0 {
		return location
	}

	u, err := url.Parse(location)
	if err != nil {
		// 解析失败，退化返回原始 location
		return location
	}

	q := u.Query()
	for k, v := range query {
		q.Set(k, v) // 覆盖策略
	}
	u.RawQuery = q.Encode()

	return u.String()
}
