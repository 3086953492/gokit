package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/ginx/pagination"
)

// OK 返回标准成功响应（HTTP 200）
//
//	response.OK(c, user)
//	response.OK(c, user, response.WithMessage("创建成功"))
func OK(c *gin.Context, data any, opts ...Option) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: o.Message,
		Data:    data,
	})
}

// OKPage 返回带分页元数据的成功响应（HTTP 200）
//
//	meta := pagination.NewMeta(p, total)
//	response.OKPage(c, list, meta)
func OKPage(c *gin.Context, data any, meta pagination.Meta, opts ...Option) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: o.Message,
		Data:    data,
		Meta:    &meta,
	})
}

