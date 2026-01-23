package i18n

// MessageData 外部传入的翻译数据
type MessageData struct {
	// Data YAML 格式的翻译内容
	Data []byte
	// Path 虚拟路径，用于识别语言和格式，如 "zh.yaml"、"en.yaml"
	// go-i18n 会从文件名中解析语言标签
	Path string
}

// TranslateConfig 翻译配置
type TranslateConfig struct {
	// MessageID 消息ID，对应 YAML 文件中的 key
	MessageID string
	// TemplateData 模板数据，用于填充消息中的变量，如 map[string]any{"Name": "张三"}
	TemplateData any
	// PluralCount 复数计数，用于确定使用哪种复数形式
	PluralCount any
	// DefaultValue 默认值，当找不到翻译时返回此值
	DefaultValue string
}
