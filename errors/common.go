package errors

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// 错误类型常量
const (
	TypeNotFound     = "NOT_FOUND"        // 404 - 记录不存在
	TypeInvalidInput = "INVALID_INPUT"    // 400 - 输入参数错误
	TypeUnauthorized = "UNAUTHORIZED"     // 401 - 未授权
	TypeForbidden    = "FORBIDDEN"        // 403 - 权限不足
	TypeDuplicate    = "DUPLICATE"        // 409 - 数据重复
	TypeInternal     = "INTERNAL_ERROR"   // 500 - 内部错误
	TypeDatabase     = "DATABASE_ERROR"   // 500 - 数据库错误
	TypeValidation   = "VALIDATION_ERROR" // 422 - 验证失败
)

// IsNotFoundError 检查是否为数据库未找到错误
func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsDuplicateError 检查是否为重复键错误
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Duplicate entry") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "UNIQUE constraint failed")
}

// FromDatabaseError 将数据库错误转换为 AppError
func FromDatabaseError(err error) *AppError {
	if err == nil {
		return nil
	}

	if IsNotFoundError(err) {
		return NotFound().Err(err).Build()
	}

	if IsDuplicateError(err) {
		return Duplicate().Err(err).Build()
	}

	return Database().Err(err).Build()
}
