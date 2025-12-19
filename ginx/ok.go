package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK 返回标准成功响应（HTTP 200）
//
//	ginx.OK(c, user)
//	// {"code":0,"message":"ok","data":{...}}
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// OKPage 返回带分页元数据的成功响应（HTTP 200）
//
//	meta := ginx.NewPageMeta(page, total)
//	ginx.OKPage(c, list, meta)
func OKPage(c *gin.Context, data any, meta PageMeta) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
		Meta:    &meta,
	})
}

