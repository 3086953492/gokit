package config

import "github.com/3086953492/gokit/config/types"

// Config 应用配置结构体，聚合各模块配置
type Config struct {
	Server     types.ServerConfig     `json:"server" yaml:"server" mapstructure:"server"`
	Database   types.DatabaseConfig   `json:"database" yaml:"database" mapstructure:"database"`
	Redis      types.RedisConfig      `json:"redis" yaml:"redis" mapstructure:"redis"`
	AuthToken  types.AuthTokenConfig  `json:"auth_token" yaml:"auth_token" mapstructure:"auth_token"`
	Log        types.LogConfig        `json:"log" yaml:"log" mapstructure:"log"`
	Middleware types.MiddlewareConfig `json:"middleware" yaml:"middleware" mapstructure:"middleware"`
	AliyunOSS  types.AliyunOSSConfig  `json:"aliyun_oss" yaml:"aliyun_oss" mapstructure:"aliyun_oss"`
	Goauth     types.GoauthConfig     `json:"goauth" yaml:"goauth" mapstructure:"goauth"`
}

// Validate 验证配置，调用子配置的 Validate 方法
func (c *Config) Validate() error {
	if err := c.Log.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultConfig 返回带有合理默认值的配置
func DefaultConfig() Config {
	return Config{
		Server: types.ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: types.DatabaseConfig{
			Port:      3306,
			Charset:   "utf8mb4",
			ParseTime: true,
			Loc:       "Local",
		},
		Redis: types.RedisConfig{
			Port: 6379,
			DB:   0,
		},
		Log: types.DefaultConfig(),
	}
}
