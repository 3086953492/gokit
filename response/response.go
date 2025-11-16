package response

import (
	"fmt"
	"net/http"
	"net/url"

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

// RedirectTemporary 临时重定向（HTTP 302）
// 可选参数 appErr 用于在重定向 URL 上附加错误信息（作为查询参数）
func RedirectTemporary(c *gin.Context, targetURL string, err error) {
	finalURL := targetURL
	if err != nil {
		var appErr *errors.AppError
		if errors.As(err, &appErr) {
			finalURL = buildRedirectURLWithError(targetURL, appErr)
		}else {
			finalURL = buildRedirectURLWithError(targetURL, errors.Internal().Msg("系统内部错误").Err(err).Build())
		}
	}
	c.Redirect(http.StatusFound, finalURL)
}

// RedirectPermanent 永久重定向（HTTP 301）
// 可选参数 appErr 用于在重定向 URL 上附加错误信息（作为查询参数）
func RedirectPermanent(c *gin.Context, targetURL string, err error) {
	finalURL := targetURL
	if err != nil {
		var appErr *errors.AppError
		if errors.As(err, &appErr) {
			finalURL = buildRedirectURLWithError(targetURL, appErr)
		}else {
			finalURL = buildRedirectURLWithError(targetURL, errors.Internal().Msg("系统内部错误").Err(err).Build())
		}
	}
	c.Redirect(http.StatusMovedPermanently, finalURL)
}

// buildRedirectURLWithError 构建带错误信息的重定向 URL
// 将 AppError 的类型和消息作为查询参数追加到目标 URL 上，并自动进行 URL 编码
func buildRedirectURLWithError(rawURL string, appErr *errors.AppError) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// 解析失败时返回原始 URL，保证容错
		return rawURL
	}

	// 获取现有查询参数
	query := parsedURL.Query()

	// 追加错误信息到查询参数
	query.Set("error", appErr.Type)
	query.Set("error_message", appErr.Message)

	// 将修改后的查询参数编码并写回 URL
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
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
