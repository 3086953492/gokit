package types

import "fmt"

// Config logger配置
type LogConfig struct {
	Level       string `json:"level" yaml:"level" mapstructure:"level"`                      // 日志级别
	Filename    string `json:"filename" yaml:"filename" mapstructure:"filename"`             // 日志文件路径
	MaxSize     int    `json:"max_size" yaml:"max_size" mapstructure:"max_size"`             // 单个文件最大大小(MB)
	MaxBackups  int    `json:"max_backups" yaml:"max_backups" mapstructure:"max_backups"`    // 最大备份文件数
	MaxAge      int    `json:"max_age" yaml:"max_age" mapstructure:"max_age"`                // 文件最大保存天数
	Compress    bool   `json:"compress" yaml:"compress" mapstructure:"compress"`             // 是否压缩
	RotateDaily bool   `json:"rotate_daily" yaml:"rotate_daily" mapstructure:"rotate_daily"` // 是否按日期轮转
	Console     bool   `json:"console" yaml:"console" mapstructure:"console"`                // 是否同时输出到控制台
	LogsDir     string `json:"logs_dir" yaml:"logs_dir" mapstructure:"logs_dir"`             // 日志目录
}

// DefaultConfig 返回默认配置
func DefaultConfig() LogConfig {
	return LogConfig{
		Level:       "info",
		Filename:    "logs/app.log",
		MaxSize:     100,
		MaxBackups:  3,
		MaxAge:      7,
		Compress:    true,
		RotateDaily: true,
		Console:     true,
		LogsDir:     "logs",
	}
}

// Validate 验证配置
func (c *LogConfig) Validate() error {
	if c.Filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if c.MaxSize <= 0 {
		c.MaxSize = 100
	}
	if c.MaxBackups < 0 {
		c.MaxBackups = 3
	}
	if c.MaxAge < 0 {
		c.MaxAge = 7
	}
	if c.LogsDir == "" {
		c.LogsDir = "logs"
	}
	return nil
}
