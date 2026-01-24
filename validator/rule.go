package validator

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Rule 自定义验证规则
type Rule struct {
	Tag     string        // 验证标签名称，如 "username_unique"
	Func    ValidatorFunc // 验证函数
	Message string        // 默认错误消息（用于未指定的语言），支持占位符 {field}, {value}, {param}
	// Messages 多语言错误消息，key 为语言代码 (BCP 47 格式)
	// 示例：map[string]string{"zh-CN": "{field}验证失败", "en": "{field} validation failed"}
	// 如果某语言未指定消息，将使用 Message 字段的值
	Messages map[string]string
}

// RegisterRule 注册单个自定义验证规则
func (m *Manager) RegisterRule(rule Rule) error {
	if rule.Tag == "" {
		return fmt.Errorf("rule tag cannot be empty")
	}
	if rule.Func == nil {
		return fmt.Errorf("validator func cannot be nil")
	}
	if rule.Message == "" {
		rule.Message = "{field}验证失败"
	}

	// 注册验证函数
	if err := m.validate.RegisterValidation(rule.Tag, rule.Func); err != nil {
		return fmt.Errorf("register validation %s: %w", rule.Tag, err)
	}

	// 为所有支持的语言注册翻译
	for _, locale := range m.locales {
		message := m.getMessageForLocale(rule, locale)
		if err := m.registerTranslationForLocale(rule.Tag, message, locale); err != nil {
			return fmt.Errorf("register translation %s for %s: %w", rule.Tag, locale, err)
		}
	}

	return nil
}

// getMessageForLocale 获取指定语言的错误消息
func (m *Manager) getMessageForLocale(rule Rule, locale string) string {
	if rule.Messages == nil {
		return rule.Message
	}

	// 尝试精确匹配
	if msg, ok := rule.Messages[locale]; ok {
		return msg
	}

	// 尝试基础语言匹配（如 zh-CN -> zh）
	baseLocale := getBaseLocale(locale)
	if msg, ok := rule.Messages[baseLocale]; ok {
		return msg
	}

	// 回退到默认消息
	return rule.Message
}

// RegisterRules 批量注册自定义验证规则
func (m *Manager) RegisterRules(rules []Rule) error {
	for _, rule := range rules {
		if err := m.RegisterRule(rule); err != nil {
			return err
		}
	}
	return nil
}

// registerTranslationForLocale 为指定语言注册验证规则的翻译
func (m *Manager) registerTranslationForLocale(tag string, message string, locale string) error {
	trans, found := m.uni.GetTranslator(locale)
	if !found {
		// 语言不存在，跳过
		return nil
	}

	return m.validate.RegisterTranslation(
		tag,
		trans,
		// 注册函数：将用户友好的占位符转换为 validator 的数字格式
		func(ut ut.Translator) error {
			// 转换占位符：{field} -> {0}, {value} -> {1}, {param} -> {2}
			msg := strings.ReplaceAll(message, "{field}", "{0}")
			msg = strings.ReplaceAll(msg, "{value}", "{1}")
			msg = strings.ReplaceAll(msg, "{param}", "{2}")
			return ut.Add(tag, msg, true)
		},
		// 翻译函数：提供占位符的实际值
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field(), fmt.Sprintf("%v", fe.Value()), fe.Param())
			return t
		},
	)
}
