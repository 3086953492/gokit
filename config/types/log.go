package types

// LogConfig 日志器配置。
type LogConfig struct {
	Level string `json:"level" yaml:"level" mapstructure:"level"` // 日志级别
}
