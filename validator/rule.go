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
		// 注册函数：添加翻译文本
		func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		},
		// 翻译函数：替换占位符
		func(ut ut.Translator, fe validator.FieldError) string {
			msg, _ := ut.T(tag)

			// 替换占位符
			msg = strings.ReplaceAll(msg, "{field}", fe.Field())
			msg = strings.ReplaceAll(msg, "{value}", fmt.Sprintf("%v", fe.Value()))
			msg = strings.ReplaceAll(msg, "{param}", fe.Param())

			return msg
		},
	)
}

