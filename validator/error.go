package validator

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ValidationResult 验证结果
type ValidationResult struct {
	errors map[string]string // 字段名 -> 中文错误消息
}

// newValidationResult 创建验证结果
func newValidationResult() *ValidationResult {
	return &ValidationResult{
		errors: make(map[string]string),
	}
}

// newValidationResultFromError 从 validator 错误创建验证结果
func newValidationResultFromError(err error, trans ut.Translator) *ValidationResult {
	result := newValidationResult()

	if err == nil {
		return result
	}

	// 转换为 ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 如果不是验证错误，返回空结果
		return result
	}

	// 翻译所有错误
	for _, e := range validationErrors {
		result.errors[e.Field()] = e.Translate(trans)
	}

	return result
}

// Valid 检查是否验证通过
func (r *ValidationResult) Valid() bool {
	return len(r.errors) == 0
}

// Errors 获取所有错误（字段名 -> 错误消息）
func (r *ValidationResult) Errors() map[string]string {
	if r.errors == nil {
		return make(map[string]string)
	}
	return r.errors
}

// First 获取第一个错误消息
func (r *ValidationResult) First() string {
	if len(r.errors) == 0 {
		return ""
	}

	// 返回第一个错误消息
	for _, msg := range r.errors {
		return msg
	}

	return ""
}

// FirstWithField 获取第一个错误消息（包含字段名）
func (r *ValidationResult) FirstWithField() string {
	if len(r.errors) == 0 {
		return ""
	}

	// 返回第一个错误消息，格式：字段: 错误
	for field, msg := range r.errors {
		return fmt.Sprintf("%s: %s", field, msg)
	}

	return ""
}

// ErrorList 获取错误列表
func (r *ValidationResult) ErrorList() []FieldError {
	if len(r.errors) == 0 {
		return []FieldError{}
	}

	errors := make([]FieldError, 0, len(r.errors))
	for field, msg := range r.errors {
		errors = append(errors, FieldError{
			Field:   field,
			Message: msg,
		})
	}

	return errors
}

// String 返回错误字符串（所有错误用分号分隔）
func (r *ValidationResult) String() string {
	if len(r.errors) == 0 {
		return ""
	}

	var messages []string
	for field, msg := range r.errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, msg))
	}

	return strings.Join(messages, "; ")
}

// Error 实现 error 接口
func (r *ValidationResult) Error() string {
	if r.Valid() {
		return ""
	}
	return r.String()
}

