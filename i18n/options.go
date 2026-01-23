package i18n

import "golang.org/x/text/language"

// Options 国际化模块配置选项
type Options struct {
	// DefaultLanguage 默认/回退语言，当找不到匹配的语言时使用
	DefaultLanguage language.Tag
	// MessageFiles YAML 翻译文件路径列表
	MessageFiles []string
	// MessageBytes 直接传入的翻译数据（从外部 []byte 加载）
	MessageBytes []MessageData
}

// Option 选项函数
type Option func(*Options)

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		DefaultLanguage: language.Chinese,
		MessageFiles:    nil,
		MessageBytes:    nil,
	}
}

// WithDefaultLanguage 设置默认/回退语言
// 当用户请求的语言在翻译文件中找不到时，将回退到此语言
func WithDefaultLanguage(lang language.Tag) Option {
	return func(o *Options) {
		o.DefaultLanguage = lang
	}
}

// WithMessageFiles 设置翻译文件路径列表
// 从文件系统加载 YAML 翻译文件，文件名需包含语言标识，如 "locales/zh.yaml"
func WithMessageFiles(paths ...string) Option {
	return func(o *Options) {
		o.MessageFiles = append(o.MessageFiles, paths...)
	}
}

// WithMessageBytes 设置直接传入的翻译数据
// 从 []byte 加载翻译，Path 字段用于识别语言和格式，如 "zh.yaml"
func WithMessageBytes(data ...MessageData) Option {
	return func(o *Options) {
		o.MessageBytes = append(o.MessageBytes, data...)
	}
}
