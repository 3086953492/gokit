package configs

type LogConfig struct {
	Level       string `mapstructure:"level"`
	Filename    string `mapstructure:"filename"`
	MaxSize     int    `mapstructure:"maxSize"`
	MaxBackups  int    `mapstructure:"maxBackups"`
	MaxAge      int    `mapstructure:"maxAge"`
	Compress    bool   `mapstructure:"compress"`
	RotateDaily bool   `mapstructure:"rotateDaily"`
}
