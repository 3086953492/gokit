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
func (v *Validator) RegisterRule(rule Rule) error {
	// 验证参数
	if rule.Tag == "" {
		return fmt.Errorf("规则标签不能为空")
	}
	if rule.Func == nil {
		return fmt.Errorf("验证函数不能为空")
	}
	if rule.Message == "" {
		rule.Message = "{field}验证失败" // 默认消息
	}

	// 注册验证函数
	if err := v.validate.RegisterValidation(rule.Tag, rule.Func); err != nil {
		return fmt.Errorf("注册验证规则 %s 失败: %w", rule.Tag, err)
	}

	// 注册翻译
	if err := v.registerTranslation(rule.Tag, rule.Message); err != nil {
		return fmt.Errorf("注册翻译 %s 失败: %w", rule.Tag, err)
	}

	return nil
}

// RegisterRules 批量注册自定义验证规则
func (v *Validator) RegisterRules(rules []Rule) error {
	for _, rule := range rules {
		if err := v.RegisterRule(rule); err != nil {
			return err
		}
	}
	return nil
}

// registerTranslation 注册验证规则的翻译
func (v *Validator) registerTranslation(tag string, message string) error {
	return v.validate.RegisterTranslation(
		tag,
		v.translator,
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
			// ut.T() 会按顺序替换 {0}, {1}, {2}
			t, _ := ut.T(tag, fe.Field(), fmt.Sprintf("%v", fe.Value()), fe.Param())
			return t
		},
	)
}

