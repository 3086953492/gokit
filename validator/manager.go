package validator

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// Manager 验证器管理器，提供统一的验证入口
type Manager struct {
	validate   *validator.Validate
	translator ut.Translator
}

// New 创建新的验证器管理器
func New(opts ...Option) (*Manager, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	v := validator.New()

	// 设置自定义字段名解析函数
	if options.TagNameFunc != nil {
		v.RegisterTagNameFunc(options.TagNameFunc)
	}

	// 初始化翻译器
	localeTranslator := getLocaleTranslator(options.Locale)
	uni := ut.New(localeTranslator, localeTranslator)
	trans, _ := uni.GetTranslator(options.Locale)

	// 注册默认翻译
	if options.RegisterDefaultTranslations {
		if err := registerDefaultTranslations(v, trans, options.Locale); err != nil {
			return nil, fmt.Errorf("register default translations: %w", err)
		}
	}

	return &Manager{
		validate:   v,
		translator: trans,
	}, nil
}

// registerDefaultTranslations 注册默认翻译
func registerDefaultTranslations(v *validator.Validate, trans ut.Translator, locale string) error {
	switch locale {
	case "en":
		return en_translations.RegisterDefaultTranslations(v, trans)
	case "zh":
		return zh_translations.RegisterDefaultTranslations(v, trans)
	default:
		return zh_translations.RegisterDefaultTranslations(v, trans)
	}
}

// Validate 验证数据并返回验证结果
func (m *Manager) Validate(data any) *Result {
	err := m.validate.Struct(data)
	return newResultFromError(err, m.translator)
}

// Struct 验证结构体（返回原始错误）
// 如果需要结构化的验证结果，请使用 Validate 方法
func (m *Manager) Struct(data any) error {
	return m.validate.Struct(data)
}

// Var 验证单个变量
func (m *Manager) Var(field any, tag string) error {
	return m.validate.Var(field, tag)
}
