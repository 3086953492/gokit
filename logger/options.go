package logger

// Options 日志管理器配置选项
type Options struct {
	// Level 日志级别，默认 InfoLevel
	Level Level
	// Console 是否输出到控制台，默认 true
	Console bool
	// File 文件输出配置，nil 表示不输出到文件
	File *FileConfig
	// AddCaller 是否记录调用位置，默认 true
	AddCaller bool
	// CallerSkip 调用栈跳过层数，默认 0
	CallerSkip int
}

// Option 配置选项函数
type Option func(*Options)

// defaultOptions 返回默认选项
func defaultOptions() *Options {
	return &Options{
		Level:      InfoLevel,
		Console:    true,
		File:       nil,
		AddCaller:  true,
		CallerSkip: 0,
	}
}

// WithLevel 设置日志级别
func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithLevelString 通过字符串设置日志级别
func WithLevelString(level string) Option {
	return func(o *Options) {
		o.Level = ParseLevel(level)
	}
}

// WithConsole 设置是否输出到控制台
func WithConsole(enabled bool) Option {
	return func(o *Options) {
		o.Console = enabled
	}
}

// WithFile 启用文件输出
func WithFile(cfg FileConfig) Option {
	return func(o *Options) {
		cfg.applyDefaults()
		o.File = &cfg
	}
}

// WithCaller 设置是否记录调用位置
func WithCaller(enabled bool) Option {
	return func(o *Options) {
		o.AddCaller = enabled
	}
}

// WithCallerSkip 设置调用栈跳过层数
func WithCallerSkip(skip int) Option {
	return func(o *Options) {
		o.CallerSkip = skip
	}
}

