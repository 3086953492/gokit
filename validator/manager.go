package validator

import (
	"fmt"

	loc "github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// Manager 验证器管理器，提供统一的验证入口
// Manager 是线程安全的，可在多个 goroutine 中共享使用
type Manager struct {
	validate      *validator.Validate
	uni           *ut.UniversalTranslator
	defaultLocale string
	locales       []string
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

	// 确定支持的语言列表
	locales := options.Locales
	if len(locales) == 0 {
		locales = []string{options.DefaultLocale}
	}

	// 确保默认语言在支持列表中
	defaultLocale := normalizeLocale(options.DefaultLocale)
	hasDefault := false
	normalizedLocales := make([]string, 0, len(locales))
	for _, loc := range locales {
		normalized := normalizeLocale(loc)
		normalizedLocales = append(normalizedLocales, normalized)
		if normalized == defaultLocale {
			hasDefault = true
		}
	}
	if !hasDefault {
		normalizedLocales = append([]string{defaultLocale}, normalizedLocales...)
	}

	// 初始化 UniversalTranslator，注册所有支持的语言
	uni, err := initUniversalTranslator(normalizedLocales, defaultLocale)
	if err != nil {
		return nil, fmt.Errorf("init universal translator: %w", err)
	}

	// 注册默认翻译
	if options.RegisterDefaultTranslations {
		for _, locale := range normalizedLocales {
			trans, found := uni.GetTranslator(locale)
			if !found {
				continue
			}
			if err := registerDefaultTranslations(v, trans, locale); err != nil {
				return nil, fmt.Errorf("register default translations for %s: %w", locale, err)
			}
		}
	}

	return &Manager{
		validate:      v,
		uni:           uni,
		defaultLocale: defaultLocale,
		locales:       normalizedLocales,
	}, nil
}

// initUniversalTranslator 初始化 UniversalTranslator，注册所有支持的语言
func initUniversalTranslator(locales []string, defaultLocale string) (*ut.UniversalTranslator, error) {
	if len(locales) == 0 {
		return nil, fmt.Errorf("at least one locale is required")
	}

	// 获取默认语言的 translator
	defaultTrans := getLocaleTranslator(defaultLocale)

	// 获取所有语言的 translators
	fallbackTranslators := make([]loc.Translator, 0, len(locales))
	for _, locale := range locales {
		trans := getLocaleTranslator(locale)
		fallbackTranslators = append(fallbackTranslators, trans)
	}

	return ut.New(defaultTrans, fallbackTranslators...), nil
}

// registerDefaultTranslations 注册默认翻译
func registerDefaultTranslations(v *validator.Validate, trans ut.Translator, locale string) error {
	baseLocale := getBaseLocale(locale)
	switch baseLocale {
	case "en":
		return en_translations.RegisterDefaultTranslations(v, trans)
	case "zh":
		return zh_translations.RegisterDefaultTranslations(v, trans)
	default:
		return zh_translations.RegisterDefaultTranslations(v, trans)
	}
}

// Validate 验证数据并返回验证结果
// 支持通过 ValidateOption 指定语言，不指定时使用默认语言
//
// 示例：
//
//	result := mgr.Validate(data) // 使用默认语言
//	result := mgr.Validate(data, WithValidateLocale("en")) // 使用英文
func (m *Manager) Validate(data any, opts ...ValidateOption) *Result {
	// 应用验证选项
	vopts := &ValidateOptions{
		Locale: m.defaultLocale,
	}
	for _, opt := range opts {
		opt(vopts)
	}

	// 获取对应语言的翻译器
	locale := normalizeLocale(vopts.Locale)
	trans, found := m.uni.GetTranslator(locale)
	if !found {
		// 回退到默认语言
		trans, _ = m.uni.GetTranslator(m.defaultLocale)
	}

	err := m.validate.Struct(data)
	return newResultFromError(err, trans)
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

// DefaultLocale 返回默认语言
func (m *Manager) DefaultLocale() string {
	return m.defaultLocale
}

// Locales 返回支持的语言列表
func (m *Manager) Locales() []string {
	result := make([]string, len(m.locales))
	copy(result, m.locales)
	return result
}
