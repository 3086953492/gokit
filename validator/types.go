package validator

import "github.com/go-playground/validator/v10"

// ValidatorFunc 验证器函数类型
type ValidatorFunc = validator.Func

// FieldLevel 字段级别接口，用于在自定义验证函数中访问字段信息
type FieldLevel = validator.FieldLevel

// FieldError 字段验证错误详情
type FieldError struct {
	Field   string      // 字段名
	Message string      // 翻译后的错误消息
	Tag     string      // 验证标签
	Value   any // 字段值
}
