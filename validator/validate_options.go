package validator

// ValidateOptions 验证时的选项
type ValidateOptions struct {
	// Locale 指定验证结果使用的语言 (BCP 47 格式)
	// 如果不指定或指定的语言不支持，将使用默认语言
	Locale string
}

// ValidateOption 验证选项函数
type ValidateOption func(*ValidateOptions)

// WithValidateLocale 设置验证结果使用的语言 (BCP 47 格式)
// 支持格式：zh-CN、en-US、zh、en 等
//
// 示例：
//
//	result := mgr.Validate(data, WithValidateLocale("en"))
//	result := mgr.Validate(data, WithValidateLocale("zh-TW"))
func WithValidateLocale(locale string) ValidateOption {
	return func(o *ValidateOptions) {
		o.Locale = locale
	}
}
