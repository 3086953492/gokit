package validator

import "strings"

// normalizeLocale 将 BCP 47 语言标签标准化为内部格式
// 支持输入格式：
//   - BCP 47 标准格式：zh-CN、en-US、zh-TW
//   - 下划线格式：zh_CN、en_US
//   - 简短格式：zh、en
//
// 返回标准化后的 locale 字符串，用于在 UniversalTranslator 中查找
func normalizeLocale(locale string) string {
	if locale == "" {
		return "zh"
	}

	// 将连字符转换为下划线（BCP 47 -> CLDR 格式）
	normalized := strings.ReplaceAll(locale, "-", "_")

	// 提取基础语言代码
	baseLocale := getBaseLocale(normalized)

	// 对于简单的语言代码（如 zh、en），直接返回基础代码
	// go-playground/locales 对于中文使用 zh，英文使用 en
	if !strings.Contains(normalized, "_") {
		return baseLocale
	}

	// 返回基础语言代码用于翻译器查找
	// go-playground/validator 的翻译包只支持 zh 和 en
	return baseLocale
}

// getBaseLocale 获取基础语言代码
// 输入: "zh-CN" -> "zh", "en-US" -> "en", "zh_TW" -> "zh"
func getBaseLocale(locale string) string {
	if locale == "" {
		return "zh"
	}

	// 处理下划线格式
	if idx := strings.Index(locale, "_"); idx > 0 {
		return locale[:idx]
	}

	// 处理连字符格式
	if idx := strings.Index(locale, "-"); idx > 0 {
		return locale[:idx]
	}

	return locale
}

// isSupportedLocale 检查语言是否被支持
// 当前支持的语言：zh（中文）、en（英文）
func isSupportedLocale(locale string) bool {
	baseLocale := getBaseLocale(locale)
	switch baseLocale {
	case "zh", "en":
		return true
	default:
		return false
	}
}
