// Package provider_gin 提供 gin 框架对 response 核心包的适配。
package provider_gin

import (
	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/response"
)

// ensure GinWriter implements response.ResponseWriter.
var _ response.ResponseWriter = (*GinWriter)(nil)

// GinWriter 将 *gin.Context 适配为 response.ResponseWriter。
type GinWriter struct {
	ctx *gin.Context
}

// NewWriter 创建 GinWriter。
func NewWriter(c *gin.Context) *GinWriter {
	return &GinWriter{ctx: c}
}

// JSON 实现 response.JSONWriter。
func (w *GinWriter) JSON(status int, body any) {
	w.ctx.JSON(status, body)
}

// Redirect 实现 response.RedirectWriter。
func (w *GinWriter) Redirect(status int, location string) {
	w.ctx.Redirect(status, location)
}

