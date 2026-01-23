package i18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Localizer 特定语言的本地化器，提供翻译方法
type Localizer struct {
	loc *i18n.Localizer
}

// T 简洁翻译方法
// messageID 为消息ID，args 为可选的 key-value 对，用于模板变量替换
// 示例：loc.T("Welcome", "Name", "张三") -> "欢迎，张三！"
// 如果翻译失败，返回 messageID 本身
func (l *Localizer) T(messageID string, args ...any) string {
	cfg := &TranslateConfig{
		MessageID: messageID,
	}

	// 将 args 转换为 map
	if len(args) >= 2 {
		data := make(map[string]any)
		for i := 0; i+1 < len(args); i += 2 {
			if key, ok := args[i].(string); ok {
				data[key] = args[i+1]
			}
		}
		if len(data) > 0 {
			cfg.TemplateData = data
		}
	}

	msg, err := l.Translate(cfg)
	if err != nil {
		return messageID
	}
	return msg
}

// Translate 完整翻译方法，支持所有翻译配置选项
func (l *Localizer) Translate(cfg *TranslateConfig) (string, error) {
	if cfg == nil || cfg.MessageID == "" {
		return "", fmt.Errorf("%w: message id is required", ErrInvalidConfig)
	}

	localizeConfig := &i18n.LocalizeConfig{
		MessageID:    cfg.MessageID,
		TemplateData: cfg.TemplateData,
		PluralCount:  cfg.PluralCount,
	}

	// 如果提供了默认值，设置 DefaultMessage
	if cfg.DefaultValue != "" {
		localizeConfig.DefaultMessage = &i18n.Message{
			ID:    cfg.MessageID,
			Other: cfg.DefaultValue,
		}
	}

	msg, err := l.loc.Localize(localizeConfig)
	if err != nil {
		// 如果有默认值且翻译失败，返回默认值
		if cfg.DefaultValue != "" {
			return cfg.DefaultValue, nil
		}
		return "", fmt.Errorf("%w: %s: %v", ErrMessageNotFound, cfg.MessageID, err)
	}

	return msg, nil
}

// MustTranslate 翻译方法，失败时 panic
func (l *Localizer) MustTranslate(cfg *TranslateConfig) string {
	msg, err := l.Translate(cfg)
	if err != nil {
		panic(err)
	}
	return msg
}

// Localizer 返回内部的 i18n.Localizer
// 一般情况下不需要直接访问，仅供特殊场景使用
func (l *Localizer) Localizer() *i18n.Localizer {
	return l.loc
}
