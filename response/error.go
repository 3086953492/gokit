package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/3086953492/gokit/errors"
)

// Error 返回错误响应
// c: gin 上下文
// err: 错误对象，支持 AppError 和普通 error
//
// 对于 AppError：
//   - 根据错误类型映射 HTTP 状态码（通过配置的 ErrorStatusMapper）
//   - 如果配置了 ShowErrorDetail=true，则在 Error 字段中包含详细信息（Message、Type、Cause、Fields）
//   - 否则只返回错误类型作为 Error 字段
//
// 对于普通 error：
//   - 返回 HTTP 500
//   - 使用配置的 FallbackErrorCode 和 FallbackErrorMessage
func Error(c *gin.Context, err error) {
	cfg := getConfig()

	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		// 处理 AppError
		statusCode := cfg.ErrorStatusMapper(appErr.Type)
		errString := buildErrorString(appErr, cfg.ShowErrorDetail)

		c.JSON(statusCode, Response{
			Success: false,
			Error:   errString,
			Message: appErr.Message,
			Data:    nil,
		})
	} else {
		// 处理普通 error
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   cfg.FallbackErrorCode,
			Message: cfg.FallbackErrorMessage,
			Data:    nil,
		})
	}
}

// buildErrorString 构建错误字符串
// 根据配置决定是否包含详细信息
func buildErrorString(appErr *errors.AppError, showDetail bool) string {
	if !showDetail {
		return appErr.Type
	}

	// 显示详细信息
	if appErr.Cause != nil {
		return fmt.Sprintf("错误消息: %s , 错误类型: %s , 错误详情: %s , 错误字段: %v",
			appErr.Message, appErr.Type, appErr.Cause.Error(), appErr.Fields)
	}

	// 没有 Cause 但有 Fields
	if len(appErr.Fields) > 0 {
		return fmt.Sprintf("错误消息: %s , 错误类型: %s , 错误字段: %v",
			appErr.Message, appErr.Type, appErr.Fields)
	}

	// 只有基本信息
	return fmt.Sprintf("错误消息: %s , 错误类型: %s", appErr.Message, appErr.Type)
}
