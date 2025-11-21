package types

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// 认证中间件配置
	Auth AuthMiddlewareConfig `json:"auth" yaml:"auth" mapstructure:"auth"`

	// CORS配置
	CORS CorsMiddlewareConfig `json:"cors" yaml:"cors" mapstructure:"cors"`
}

// AuthMiddlewareConfig 认证中间件配置
type AuthMiddlewareConfig struct {
	SkipPaths []string `json:"skip_paths" yaml:"skip_paths" mapstructure:"skip_paths"`
}

type CorsMiddlewareConfig struct {
	AllowOrigins []string `json:"allow_origins" yaml:"allow_origins" mapstructure:"allow_origins"`
	AllowMethods []string `json:"allow_methods" yaml:"allow_methods" mapstructure:"allow_methods"`
	AllowHeaders []string `json:"allow_headers" yaml:"allow_headers" mapstructure:"allow_headers"`
}
