package provider_gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/response"
)

// Responder 封装 Manager + GinWriter，提供便捷方法。
type Responder struct {
	mgr *response.Manager
}

// NewResponder 创建 Responder。
func NewResponder(mgr *response.Manager) *Responder {
	return &Responder{mgr: mgr}
}

// --------------------------------------------------------------------
// Problem / Error 响应
// --------------------------------------------------------------------

// Error 写入 RFC7807 problem+json 响应。
// instance 通常为当前请求路径。
func (r *Responder) Error(c *gin.Context, err error) {
	r.ErrorWithInstance(c, err, c.Request.URL.Path)
}

// ErrorWithInstance 写入 RFC7807 problem+json 响应，自定义 instance。
func (r *Responder) ErrorWithInstance(c *gin.Context, err error, instance string) {
	// 设置 Content-Type 为 problem+json
	c.Header("Content-Type", "application/problem+json")
	r.mgr.WriteProblem(NewWriter(c), err, instance)
}

// --------------------------------------------------------------------
// 成功响应
// --------------------------------------------------------------------

// Success 写入成功响应（200）。
func (r *Responder) Success(c *gin.Context, message string, data any) {
	r.mgr.WriteSuccess(NewWriter(c), message, data)
}

// Paginated 写入分页成功响应（200）。
func (r *Responder) Paginated(c *gin.Context, items any, total int64, page, pageSize int) {
	r.mgr.WritePaginated(NewWriter(c), items, total, page, pageSize)
}

// --------------------------------------------------------------------
// 重定向
// --------------------------------------------------------------------

// RedirectTemporary 临时重定向（302），可选 err 用于注入 OAuth2 参数。
func (r *Responder) RedirectTemporary(c *gin.Context, targetURL string, err error, params map[string]string) {
	r.mgr.WriteRedirectTemporary(NewWriter(c), targetURL, err, params)
}

// RedirectPermanent 永久重定向（301），可选 err 用于注入 OAuth2 参数。
func (r *Responder) RedirectPermanent(c *gin.Context, targetURL string, err error, params map[string]string) {
	r.mgr.WriteRedirectPermanent(NewWriter(c), targetURL, err, params)
}

// Redirect 自定义状态码重定向。
func (r *Responder) Redirect(c *gin.Context, status int, targetURL string, err error, params map[string]string) {
	r.mgr.WriteRedirect(NewWriter(c), status, targetURL, err, params)
}

// --------------------------------------------------------------------
// 快捷包级函数（使用默认 Manager）
// --------------------------------------------------------------------

var defaultResponder = NewResponder(response.NewManager())

// Error 使用默认配置写入错误响应。
func Error(c *gin.Context, err error) {
	defaultResponder.Error(c, err)
}

// Success 使用默认配置写入成功响应。
func Success(c *gin.Context, message string, data any) {
	defaultResponder.Success(c, message, data)
}

// Paginated 使用默认配置写入分页响应。
func Paginated(c *gin.Context, items any, total int64, page, pageSize int) {
	defaultResponder.Paginated(c, items, total, page, pageSize)
}

// RedirectTemporary 使用默认配置临时重定向（302）。
func RedirectTemporary(c *gin.Context, targetURL string, err error, params map[string]string) {
	defaultResponder.RedirectTemporary(c, targetURL, err, params)
}

// RedirectPermanent 使用默认配置永久重定向（301）。
func RedirectPermanent(c *gin.Context, targetURL string, err error, params map[string]string) {
	defaultResponder.RedirectPermanent(c, targetURL, err, params)
}

// Redirect 使用默认配置自定义状态码重定向。
func Redirect(c *gin.Context, status int, targetURL string, err error, params map[string]string) {
	defaultResponder.Redirect(c, status, targetURL, err, params)
}

// SetDefaultManager 替换默认 Manager（用于自定义配置后全局替换）。
func SetDefaultManager(mgr *response.Manager) {
	defaultResponder = NewResponder(mgr)
}

// HTTP 状态常量（方便外部直接调用 Redirect 时使用）。
const (
	StatusMovedPermanently  = http.StatusMovedPermanently
	StatusFound             = http.StatusFound
	StatusTemporaryRedirect = http.StatusTemporaryRedirect
	StatusPermanentRedirect = http.StatusPermanentRedirect
)

