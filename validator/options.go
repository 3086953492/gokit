package validator

import (
	"reflect"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
)

// Options 验证器配置选项
type Options struct {
	// Locale 语言/地区标识，默认 "zh"
	Locale string
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
		Locale:                      "zh",
		RegisterDefaultTranslations: true,
		TagNameFunc:                 nil,
	}
}

// WithLocale 设置语言/地区
func WithLocale(locale string) Option {
	return func(o *Options) {
		o.Locale = locale
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
func getLocaleTranslator(locale string) locales.Translator {
	switch locale {
	case "en":
		return en.New()
	case "zh":
		return zh.New()
	default:
		return zh.New()
	}
}
