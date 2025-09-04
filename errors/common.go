package errors

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// 通用错误类型常量
const (
	TypeNotFound       = "NOT_FOUND"
	TypeDuplicateKey   = "DUPLICATE_KEY"
	TypeDatabaseError  = "DATABASE_ERROR"
	TypeInternalError  = "INTERNAL_ERROR"
	TypeValidation     = "VALIDATION_ERROR"
	TypeInvalidInput   = "INVALID_INPUT"
	TypeAutoMigrate    = "AUTO_MIGRATE"
	TypeServerInternal = "SERVER_INTERNAL"
)

// 通用错误实例
var (
	ErrNotFound       = New(TypeNotFound, "记录不存在")
	ErrDuplicateKey   = New(TypeDuplicateKey, "数据已存在")
	ErrDatabaseError  = New(TypeDatabaseError, "数据库操作失败")
	ErrInternalError  = New(TypeInternalError, "系统内部错误")
	ErrValidation     = New(TypeValidation, "数据验证失败")
	ErrInvalidInput   = New(TypeInvalidInput, "输入参数错误")
	ErrAutoMigrate    = New(TypeAutoMigrate, "自动迁移失败")
	ErrServerInternal = New(TypeServerInternal, "服务器内部错误")
)

// 错误类型检查函数
func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Duplicate entry") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "UNIQUE constraint failed")
}

// 数据库错误转换
func FromDatabaseError(err error) *AppError {
	if err == nil {
		return nil
	}

	if IsNotFoundError(err) {
		return ErrNotFound
	}

	if IsDuplicateError(err) {
		return ErrDuplicateKey
	}

	return Wrap(err, TypeDatabaseError, "数据库操作失败")
}
