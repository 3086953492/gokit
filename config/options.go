package config

// Logger 日志接口，用于配置包内部日志输出
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// Options 配置管理器选项
type Options struct {
	// ConfigFile 显式指定配置文件路径（最高优先级）
	ConfigFile string

	// ConfigDir 配置文件目录
	ConfigDir string

	// EnvConfigKey 从环境变量读取配置文件路径的 key（默认 "CONFIG"）
	EnvConfigKey string

	// Formats 支持的配置文件格式（默认 ["yaml", "json"]）
	Formats []string

	// Logger 日志输出接口（默认 no-op）
	Logger Logger
}

// Option 配置选项函数
type Option func(*Options)

// DefaultOptions 返回默认选项
func DefaultOptions() *Options {
	return &Options{
		ConfigDir:    ".",
		EnvConfigKey: "CONFIG",
		Formats:      []string{"yaml", "json"},
		Logger:       noopLogger{},
	}
}

// WithConfigFile 设置显式配置文件路径（最高优先级）
func WithConfigFile(path string) Option {
	return func(o *Options) {
		o.ConfigFile = path
	}
}

// WithConfigDir 设置配置文件目录
func WithConfigDir(dir string) Option {
	return func(o *Options) {
		o.ConfigDir = dir
	}
}

// WithEnvConfigKey 设置从环境变量读取配置文件路径的 key
func WithEnvConfigKey(key string) Option {
	return func(o *Options) {
		o.EnvConfigKey = key
	}
}

// WithFormats 设置支持的配置文件格式
func WithFormats(formats []string) Option {
	return func(o *Options) {
		o.Formats = formats
	}
}

// WithLogger 设置日志输出接口
func WithLogger(logger Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// noopLogger 空日志实现
type noopLogger struct{}

func (noopLogger) Info(msg string, args ...any)  {}
func (noopLogger) Error(msg string, args ...any) {}
