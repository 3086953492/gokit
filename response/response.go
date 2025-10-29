package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/errors"
)

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 返回成功响应
func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, err error) {

	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		// 根据错误类型返回相应的 HTTP 状态码
		statusCode := getHTTPStatus(appErr.Type)

		var errString string

		// 在开发环境返回详细信息
		if config.GetGlobalConfig().Server.Mode == gin.DebugMode && appErr.Cause != nil {
			errString = fmt.Sprintf("错误消息: %s , 错误类型: %s , 错误详情: %s , 错误字段: %v", appErr.Message, appErr.Type, appErr.Cause.Error(), appErr.Fields)
		} else {
			errString = appErr.Type
		}

		c.JSON(statusCode, Response{
			Success: false,
			Error:   errString,
			Message: appErr.Message,
			Data:    nil,
		})
	} else {
		c.JSON(500, Response{
			Success: false,
			Error:   "SYSTEM_ERROR",
			Message: "系统内部错误",
			Data:    nil,
		})
	}
}

// Paginated 返回分页响应
func Paginated(c *gin.Context, data any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "获取成功",
		Data: gin.H{
			"items":       data,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": total / int64(pageSize),
		},
	})
}

// getHTTPStatus 将错误类型映射为HTTP状态码
func getHTTPStatus(errType string) int {
	switch errType {
	case errors.TypeNotFound:
		return 404
	case errors.TypeInvalidInput:
		return 400
	case errors.TypeUnauthorized:
		return 401
	case errors.TypeForbidden:
		return 403
	case errors.TypeDuplicate:
		return 409
	case errors.TypeValidation:
		return 422
	default:
		return 500
	}
}
