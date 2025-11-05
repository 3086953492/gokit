package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ValidationResult 验证结果
type ValidationResult struct {
	Valid   bool   // 是否验证通过
	Message string // 友好的错误消息（给用户看）
	Err     error  // 原始错误（方便包装和日志记录）
}

// newValidationResult 创建空的验证结果（验证通过）
func newValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid: true,
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

	// 获取第一个错误的翻译消息
	var firstMessage string
	for _, e := range validationErrors {
		firstMessage = e.Translate(trans)
		break // 只取第一个
	}

	return &ValidationResult{
		Valid:   false,
		Message: firstMessage,
		Err:     err,
	}
}

// ErrorList 获取错误列表（从原始错误中提取）
func (r *ValidationResult) ErrorList() []FieldError {
	if r.Err == nil {
		return []FieldError{}
	}

	// 尝试转换为 ValidationErrors
	validationErrors, ok := r.Err.(validator.ValidationErrors)
	if !ok {
		return []FieldError{}
	}

	// 获取翻译器（使用默认验证器的翻译器）
	trans := GetDefaultValidator().GetTranslator()

	errors := make([]FieldError, 0, len(validationErrors))
	for _, e := range validationErrors {
		errors = append(errors, FieldError{
			Field:   e.Field(),
			Message: e.Translate(trans),
			Tag:     e.Tag(),
			Value:   e.Value(),
		})
	}

	return errors
}

// String 返回错误字符串
func (r *ValidationResult) String() string {
	if r.Valid {
		return ""
	}
	return r.Message
}

// Error 实现 error 接口
func (r *ValidationResult) Error() string {
	return r.String()
}
