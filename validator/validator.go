package validator

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// Validator 验证器实例
type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

// New 创建新的验证器实例
func New() *Validator {
	v := validator.New()

	// 初始化中文翻译器
	zhLocale := zh.New()
	uni := ut.New(zhLocale, zhLocale)
	trans, _ := uni.GetTranslator("zh")

	// 注册默认中文翻译
	_ = zh_translations.RegisterDefaultTranslations(v, trans)

	return &Validator{
		validate:   v,
		translator: trans,
	}
}

// Validate 验证数据并返回验证结果
func (v *Validator) Validate(data interface{}) *ValidationResult {
	err := v.validate.Struct(data)
	return newValidationResultFromError(err, v.translator)
}

// Struct 验证结构体（返回原始错误）
// 如果需要结构化的验证结果，请使用 Validate 方法
func (v *Validator) Struct(data interface{}) error {
	return v.validate.Struct(data)
}

// Var 验证单个变量
func (v *Validator) Var(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// GetTranslator 获取翻译器实例
func (v *Validator) GetTranslator() ut.Translator {
	return v.translator
}

// GetValidate 获取底层的 validator.Validate 实例
// 用于需要访问原始验证器的高级场景
func (v *Validator) GetValidate() *validator.Validate {
	return v.validate
}

