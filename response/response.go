package response

import (
	"errors"
	apperrors "github.com/3086953492/YaBase/errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var appErr *apperrors.AppError

	if errors.As(err, &appErr) {
		// 应用错误：使用错误类型和消息
		httpStatus := getHTTPStatus(appErr.Type)
		c.JSON(httpStatus, Response{
			Success: false,
			Error:   appErr.Type,
			Message: appErr.Message,
		})
	} else {
		// 未知错误：返回通用错误
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   apperrors.TypeInternalError,
			Message: "系统内部错误",
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
func getHTTPStatus(errorType string) int {
	switch errorType {
	case apperrors.TypeUnauthorized, apperrors.TypeTokenExpired, apperrors.TypeTokenInvalid:
		return http.StatusUnauthorized // 401
	case apperrors.TypePermissionDenied:
		return http.StatusForbidden // 403
	case apperrors.TypeNotFound, apperrors.TypeUserNotFound:
		return http.StatusNotFound // 404
	case apperrors.TypeInvalidInput, apperrors.TypeValidation, apperrors.TypeDuplicateKey:
		return http.StatusBadRequest // 400
	case apperrors.TypeInternalError, apperrors.TypeDatabaseError:
		return http.StatusInternalServerError // 500
	default:
		return http.StatusBadRequest // 400
	}
}
