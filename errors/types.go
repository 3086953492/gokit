package errors

import (
	"errors"
	"fmt"
)

// AppError 应用错误类型
type AppError struct {
	Type    string         `json:"type"`             // 错误类型
	Message string         `json:"message"`          // 错误消息
	Cause   error          `json:"-"`                // 原始错误，不序列化
	Fields  map[string]any `json:"fields,omitempty"` // 上下文字段
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap 实现错误解包，支持 errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Cause
}

// GetField 获取指定的上下文字段
func (e *AppError) GetField(key string) (any, bool) {
	if e.Fields == nil {
		return nil, false
	}
	val, ok := e.Fields[key]
	return val, ok
}

// HasField 检查是否存在指定的上下文字段
func (e *AppError) HasField(key string) bool {
	_, ok := e.GetField(key)
	return ok
}

// IsAppError 检查错误类型是否为指定的 AppError
func IsAppError(err error, target *AppError) bool {
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

// GetFields 获取错误的所有上下文字段
func GetFields(err error) map[string]any {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Fields
	}
	return nil
}
