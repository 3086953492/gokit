package validator

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ValidationResult 验证结果
type ValidationResult struct {
	Valid   bool              // 是否验证通过
	Message string            // 第一个错误消息（快捷访问）
	Errors  map[string]string // 所有错误（字段名 -> 错误消息）
}

// newValidationResult 创建空的验证结果（验证通过）
func newValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Errors: make(map[string]string),
	}
}

// newValidationResultFromError 从 validator 错误创建验证结果
func newValidationResultFromError(err error, trans ut.Translator) *ValidationResult {
	if err == nil {
		return newValidationResult()
	}

	// 转换为 ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 如果不是验证错误，返回验证通过的结果
		return newValidationResult()
	}

	// 翻译所有错误
	errors := make(map[string]string)
	var firstMessage string
	for _, e := range validationErrors {
		msg := e.Translate(trans)
		errors[e.Field()] = msg
		// 记录第一个错误消息
		if firstMessage == "" {
			firstMessage = msg
		}
	}

	return &ValidationResult{
		Valid:   len(errors) == 0,
		Message: firstMessage,
		Errors:  errors,
	}
}

// ErrorList 获取错误列表
func (r *ValidationResult) ErrorList() []FieldError {
	if len(r.Errors) == 0 {
		return []FieldError{}
	}

	errors := make([]FieldError, 0, len(r.Errors))
	for field, msg := range r.Errors {
		errors = append(errors, FieldError{
			Field:   field,
			Message: msg,
		})
	}

	return errors
}

// String 返回错误字符串（所有错误用分号分隔）
func (r *ValidationResult) String() string {
	if len(r.Errors) == 0 {
		return ""
	}

	var messages []string
	for field, msg := range r.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, msg))
	}

	return strings.Join(messages, "; ")
}

// Error 实现 error 接口
func (r *ValidationResult) Error() string {
	if r.Valid {
		return ""
	}
	return r.String()
}
