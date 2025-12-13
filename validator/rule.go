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
	Message string        // 中文错误消息，支持占位符 {field}, {value}, {param}
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

	// 注册翻译
	if err := m.registerTranslation(rule.Tag, rule.Message); err != nil {
		return fmt.Errorf("register translation %s: %w", rule.Tag, err)
	}

	return nil
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

// registerTranslation 注册验证规则的翻译
func (m *Manager) registerTranslation(tag string, message string) error {
	return m.validate.RegisterTranslation(
		tag,
		m.translator,
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
