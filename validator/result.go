package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Result 验证结果
type Result struct {
	Valid   bool         // 是否验证通过
	Message string       // 友好的错误消息（给用户看）
	Err     error        // 原始错误（方便包装和日志记录）
	errors  []FieldError // 已翻译的字段错误列表
}

// newResult 创建空的验证结果（验证通过）
func newResult() *Result {
	return &Result{
		Valid:  true,
		errors: []FieldError{},
	}
}

// newResultFromError 从 validator 错误创建验证结果
func newResultFromError(err error, trans ut.Translator) *Result {
	if err == nil {
		return newResult()
	}

	// 转换为 ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非 ValidationErrors 也应标记为失败并保留原始错误
		return &Result{
			Valid:   false,
			Message: err.Error(),
			Err:     err,
			errors:  []FieldError{},
		}
	}

	// 预先翻译所有字段错误
	fieldErrors := make([]FieldError, 0, len(validationErrors))
	var firstMessage string
	for i, e := range validationErrors {
		msg := e.Translate(trans)
		if i == 0 {
			firstMessage = msg
		}
		fieldErrors = append(fieldErrors, FieldError{
			Field:   e.Field(),
			Message: msg,
			Tag:     e.Tag(),
			Value:   e.Value(),
		})
	}

	return &Result{
		Valid:   false,
		Message: firstMessage,
		Err:     err,
		errors:  fieldErrors,
	}
}

// ErrorList 获取错误列表
func (r *Result) ErrorList() []FieldError {
	if r.errors == nil {
		return []FieldError{}
	}
	return r.errors
}

// String 返回错误字符串
func (r *Result) String() string {
	if r.Valid {
		return ""
	}
	return r.Message
}

// Error 实现 error 接口
func (r *Result) Error() string {
	return r.String()
}
