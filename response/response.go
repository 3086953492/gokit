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
// targetURL: 目标重定向地址
// err: 可选的错误对象，若非 nil 则会在 URL 中添加 error 和 error_description 参数
// params: 自定义查询参数（如 map[string]string{"code":"123","state":"ok"}），会覆盖原 URL 中的同名参数
//
// 注意：若 URL 原本带有同名参数，将被新传入的参数覆盖；错误相关参数始终以最新计算值为准
// 若暂不需要额外参数，可传入 nil 或 map[string]string{}
func RedirectTemporary(c *gin.Context, targetURL string, err error, params map[string]string) {
	var appErr *errors.AppError
	if err != nil {
		if !errors.As(err, &appErr) {
			// 若 err 不是 AppError，则包装为内部错误
			appErr = errors.Internal().Msg("系统内部错误").Err(err).Build()
		}
	}
	finalURL := buildRedirectURL(targetURL, appErr, params)
	c.Redirect(http.StatusFound, finalURL)
}

// RedirectPermanent 永久重定向（HTTP 301）
// targetURL: 目标重定向地址
// err: 可选的错误对象，若非 nil 则会在 URL 中添加 error 和 error_description 参数
// params: 自定义查询参数（如 map[string]string{"code":"123","state":"ok"}），会覆盖原 URL 中的同名参数
//
// 注意：若 URL 原本带有同名参数，将被新传入的参数覆盖；错误相关参数始终以最新计算值为准
// 若暂不需要额外参数，可传入 nil 或 map[string]string{}
func RedirectPermanent(c *gin.Context, targetURL string, err error, params map[string]string) {
	var appErr *errors.AppError
	if err != nil {
		if !errors.As(err, &appErr) {
			// 若 err 不是 AppError，则包装为内部错误
			appErr = errors.Internal().Msg("系统内部错误").Err(err).Build()
		}
	}
	finalURL := buildRedirectURL(targetURL, appErr, params)
	c.Redirect(http.StatusMovedPermanently, finalURL)
}

// buildRedirectURL 构建带参数的重定向 URL
// 支持添加自定义参数和错误信息到目标 URL 的查询参数中，并自动进行 URL 编码
// params: 自定义查询参数，会覆盖原 URL 中的同名参数
// appErr: 可选的错误对象，若非 nil 则会添加 error 和 error_description 参数
func buildRedirectURL(rawURL string, appErr *errors.AppError, params map[string]string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// 解析失败时返回原始 URL，保证容错
		return rawURL
	}

	// 获取现有查询参数
	query := parsedURL.Query()

	// 添加自定义参数（会覆盖原有同名参数）
	for k, v := range params {
		query.Set(k, v)
	}

	// 若存在错误信息，添加错误相关参数（优先级最高，会覆盖 params 中的同名参数）
	if appErr != nil {
		errorCode := mapToOAuthErrorCode(appErr.Type)
		query.Set("error", errorCode)
		query.Set("error_description", appErr.Message)
	}

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

// mapToOAuthErrorCode 将 AppError 类型映射为 OAuth 2.0 标准错误代码
func mapToOAuthErrorCode(errType string) string {
	switch errType {
	case errors.TypeInvalidInput:
		return "invalid_request"
	case errors.TypeUnauthorized:
		return "unauthorized_client"
	case errors.TypeForbidden:
		return "access_denied"
	default:
		return "server_error"
	}
}
