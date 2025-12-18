// Package logger 提供统一的日志管理功能
package logger

// Level 日志级别
type Level int8

const (
	// DebugLevel 调试级别，输出所有日志
	DebugLevel Level = iota - 1
	// InfoLevel 信息级别（默认）
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
)

// String 返回级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "info"
	}
}

// ParseLevel 从字符串解析日志级别
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// RotateStrategy 日志轮转策略
type RotateStrategy int

const (
	// RotateBySize 按大小轮转（默认）
	RotateBySize RotateStrategy = iota
	// RotateByDate 按日期轮转（每日一个文件）
	RotateByDate
)

// FileConfig 文件输出配置
type FileConfig struct {
	// Filename 日志文件路径（必填）
	Filename string
	// MaxSize 单个文件最大大小（MB），默认 100
	MaxSize int
	// MaxBackups 最大备份文件数，默认 3
	MaxBackups int
	// MaxAge 文件最大保存天数，默认 7
	MaxAge int
	// Compress 是否压缩旧文件，默认 true
	Compress bool
	// RotateStrategy 轮转策略，默认按大小
	RotateStrategy RotateStrategy
}

// applyDefaults 应用默认值
func (c *FileConfig) applyDefaults() {
	if c.MaxSize <= 0 {
		c.MaxSize = 100
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 3
	}
	if c.MaxAge <= 0 {
		c.MaxAge = 7
	}
}

