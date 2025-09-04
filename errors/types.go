package errors

import (
	"errors"
	"fmt"
)

// AppError 应用错误类型
type AppError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Cause   error  `json:"-"` // 原始错误，不序列化到JSON
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// New 创建新的应用错误
func New(errorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
	}
}

// Wrap 包装已有错误
func Wrap(err error, errorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Cause:   err,
	}
}

// Is 检查错误类型
func Is(err error, target *AppError) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == target.Type
	}
	return false
}

// GetType 获取错误类型
func GetType(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return "UNKNOWN"
}
