package redirect

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Redirect 执行 HTTP 重定向并终止后续 handler 处理。
//
// 默认使用 302 Found，可通过 WithStatus 覆盖。
// 可通过 WithQuery 附加 query 参数。
//
//	redirect.Redirect(c, "/login")
//	redirect.Redirect(c, "/dashboard", redirect.WithStatus(http.StatusMovedPermanently))
//	redirect.Redirect(c, "/callback", redirect.WithQuery(map[string]string{"token": "abc"}))
func Redirect(c *gin.Context, location string, opts ...Option) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	finalLocation := buildLocation(location, o.Query)
	c.Redirect(o.Status, finalLocation)
	c.Abort()
}

// To 执行指定状态码的 HTTP 重定向并终止后续 handler 处理。
//
//	redirect.To(c, http.StatusMovedPermanently, "/new-path")
func To(c *gin.Context, status int, location string) {
	c.Redirect(status, location)
	c.Abort()
}

// ToWithQuery 执行指定状态码的 HTTP 重定向并附加 query 参数。
//
//	redirect.ToWithQuery(c, 302, "/callback", map[string]string{"code": "123"})
func ToWithQuery(c *gin.Context, status int, location string, query map[string]string) {
	finalLocation := buildLocation(location, query)
	c.Redirect(status, finalLocation)
	c.Abort()
}

// ---------------------------------------------------------------------------
// 便捷方法
// ---------------------------------------------------------------------------

// Found 执行 302 Found 重定向
func Found(c *gin.Context, location string) {
	To(c, http.StatusFound, location)
}

// Permanent 执行 301 Moved Permanently 重定向
func Permanent(c *gin.Context, location string) {
	To(c, http.StatusMovedPermanently, location)
}

// Temporary 执行 307 Temporary Redirect 重定向（保持请求方法和请求体）
func Temporary(c *gin.Context, location string) {
	To(c, http.StatusTemporaryRedirect, location)
}

// SeeOther 执行 303 See Other 重定向（常用于 POST 后跳转到 GET）
func SeeOther(c *gin.Context, location string) {
	To(c, http.StatusSeeOther, location)
}

// ---------------------------------------------------------------------------
// 内部工具函数
// ---------------------------------------------------------------------------

// buildLocation 解析 location 并合并 query 参数（同名 key 覆盖）
func buildLocation(location string, query map[string]string) string {
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

