package configs

import "time"

// 中间件配置 - 先简单点，后面再扩展
type MiddlewareConfig struct {
	// JWT配置
	JWT struct {
		Secret    string        `yaml:"secret"`
		Expire    time.Duration `yaml:"expire"`
		SkipPaths []string      `yaml:"skip_paths"`
	} `yaml:"jwt"`

	// CORS配置
	CORS struct {
		AllowOrigins []string `yaml:"allow_origins"`
		AllowMethods []string `yaml:"allow_methods"`
		AllowHeaders []string `yaml:"allow_headers"`
	} `yaml:"cors"`
}

// 默认配置
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		JWT: struct {
			Secret    string        `yaml:"secret"`
			Expire    time.Duration `yaml:"expire"`
			SkipPaths []string      `yaml:"skip_paths"`
		}{
			Secret: "your-default-secret",
			Expire: 24 * time.Hour,
			SkipPaths: []string{
				"/api/account/v1/auth/login",
				"/api/account/v1/auth/register",
			},
		},
		CORS: struct {
			AllowOrigins []string `yaml:"allow_origins"`
			AllowMethods []string `yaml:"allow_methods"`
			AllowHeaders []string `yaml:"allow_headers"`
		}{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
		},
	}
}
