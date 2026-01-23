package i18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/yaml.v3"
)

// Manager 国际化管理器，持有翻译包 Bundle，提供 Localizer 工厂方法。
// Manager 是线程安全的，可在多个 goroutine 中共享使用。
type Manager struct {
	bundle *i18n.Bundle
	opts   *Options
}

// New 创建国际化管理器
// 至少需要通过 WithMessageFiles 或 WithMessageBytes 提供翻译数据
func New(opts ...Option) (*Manager, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 创建 Bundle，设置默认语言
	bundle := i18n.NewBundle(options.DefaultLanguage)

	// 注册 YAML 解析函数
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	// 从文件路径加载翻译
	for _, path := range options.MessageFiles {
		if _, err := bundle.LoadMessageFile(path); err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrLoadFailed, path, err)
		}
	}

	// 从 []byte 加载翻译
	for _, msg := range options.MessageBytes {
		if _, err := bundle.ParseMessageFileBytes(msg.Data, msg.Path); err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrLoadFailed, msg.Path, err)
		}
	}

	return &Manager{
		bundle: bundle,
		opts:   options,
	}, nil
}

// Localizer 根据语言偏好创建 Localizer
// langs 可以是语言标签（如 "zh-CN", "en"）或 Accept-Language 头的值（如 "zh-CN,en;q=0.9"）
// go-i18n 会自动解析并选择最匹配的语言
func (m *Manager) Localizer(langs ...string) *Localizer {
	loc := i18n.NewLocalizer(m.bundle, langs...)
	return &Localizer{loc: loc}
}

// Bundle 返回内部的 i18n.Bundle
// 一般情况下不需要直接访问，仅供特殊场景使用
func (m *Manager) Bundle() *i18n.Bundle {
	return m.bundle
}
