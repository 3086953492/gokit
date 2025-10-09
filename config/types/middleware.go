package types

import "time"

// 中间件配置 - 先简单点，后面再扩展
type MiddlewareConfig struct {
	// JWT配置
	JWT struct {
		Secret    string        `json:"secret" yaml:"secret" mapstructure:"secret"`
		Expire    time.Duration `json:"expire" yaml:"expire" mapstructure:"expire"`
		SkipPaths []string      `json:"skip_paths" yaml:"skip_paths" mapstructure:"skip_paths"`
	} `json:"jwt" yaml:"jwt" mapstructure:"jwt"`

	// CORS配置
	CORS struct {
		AllowOrigins []string `json:"allow_origins" yaml:"allow_origins" mapstructure:"allow_origins"`
		AllowMethods []string `json:"allow_methods" yaml:"allow_methods" mapstructure:"allow_methods"`
		AllowHeaders []string `json:"allow_headers" yaml:"allow_headers" mapstructure:"allow_headers"`
	} `json:"cors" yaml:"cors" mapstructure:"cors"`
}
