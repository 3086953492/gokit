package validator

import (
	"reflect"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
)

// Options 验证器配置选项
type Options struct {
	// DefaultLocale 默认语言/地区标识 (BCP 47 格式，如 "zh-CN"、"en")，默认 "zh-CN"
	DefaultLocale string
	// Locales 支持的语言列表 (BCP 47 格式)，如 []string{"zh-CN", "en"}
	// 如果为空，将自动使用 DefaultLocale 作为唯一支持的语言
	Locales []string
	// RegisterDefaultTranslations 是否注册默认翻译，默认 true
	RegisterDefaultTranslations bool
	// TagNameFunc 自定义字段名解析函数（如从 json tag 获取字段名）
	TagNameFunc func(fld reflect.StructField) string
}

// Option 选项函数
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		DefaultLocale:               "zh-CN",
		Locales:                     nil,
		RegisterDefaultTranslations: true,
		TagNameFunc:                 nil,
	}
}

// WithDefaultLocale 设置默认语言/地区 (BCP 47 格式)
// 支持格式：zh-CN、en-US、zh、en 等
func WithDefaultLocale(locale string) Option {
	return func(o *Options) {
		o.DefaultLocale = locale
	}
}

// WithLocale 设置默认语言/地区（WithDefaultLocale 的别名，保持向后兼容）
// Deprecated: 推荐使用 WithDefaultLocale
func WithLocale(locale string) Option {
	return WithDefaultLocale(locale)
}

// WithLocales 设置支持的语言列表 (BCP 47 格式)
// 示例：WithLocales("zh-CN", "en", "zh-TW")
func WithLocales(locales ...string) Option {
	return func(o *Options) {
		o.Locales = locales
	}
}

// WithRegisterDefaultTranslations 设置是否注册默认翻译
func WithRegisterDefaultTranslations(register bool) Option {
	return func(o *Options) {
		o.RegisterDefaultTranslations = register
	}
}

// WithTagNameFunc 设置字段名解析函数
// 常用于从 json/form tag 获取字段名，例如：
//
//	WithTagNameFunc(func(fld reflect.StructField) string {
//	    name := fld.Tag.Get("json")
//	    if name == "" || name == "-" {
//	        return fld.Name
//	    }
//	    return strings.Split(name, ",")[0]
//	})
func WithTagNameFunc(fn func(fld reflect.StructField) string) Option {
	return func(o *Options) {
		o.TagNameFunc = fn
	}
}

// getLocaleTranslator 根据 locale 获取对应的 locales.Translator
// 支持 BCP 47 格式的语言标签
func getLocaleTranslator(locale string) locales.Translator {
	baseLocale := getBaseLocale(locale)
	switch baseLocale {
	case "en":
		return en.New()
	case "zh":
		return zh.New()
	default:
		return zh.New()
	}
}
