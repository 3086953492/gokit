package validator

import "strings"

// WithPrefix 设置标签前缀
func WithPrefix(prefix string) RegisterOption {
	return func(config *PackageConfig) {
		config.Prefix = prefix
	}
}

// WithCustomTags 设置自定义标签映射
func WithCustomTags(tags map[string]string) RegisterOption {
	return func(config *PackageConfig) {
		config.Tags = tags
	}
}

// WithSkip 设置跳过的方法
func WithSkip(methods ...string) RegisterOption {
	return func(config *PackageConfig) {
		config.Skip = methods
	}
}

// WithTransformer 设置自定义标签转换器
func WithTransformer(transformer TagTransformer) RegisterOption {
	return func(config *PackageConfig) {
		config.Transform = transformer
	}
}

// 预定义的转换器
var (
	// CamelCaseTransformer 驼峰转换器（默认）
	CamelCaseTransformer TagTransformer = CamelToSnake

	// LowerCaseTransformer 小写转换器
	LowerCaseTransformer TagTransformer = strings.ToLower

	// IdentityTransformer 恒等转换器（不转换）
	IdentityTransformer TagTransformer = func(s string) string { return s }
)
