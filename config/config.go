package config

import "github.com/3086953492/gokit/config/types"

type Config struct {
	Server     types.ServerConfig     `json:"server" yaml:"server" mapstructure:"server"`
	Database   types.DatabaseConfig   `json:"database" yaml:"database" mapstructure:"database"`
	Redis      types.RedisConfig      `json:"redis" yaml:"redis" mapstructure:"redis"`
	AuthToken  types.AuthTokenConfig  `json:"auth_token" yaml:"auth_token" mapstructure:"auth_token"`
	OAuth      types.OAuthConfig      `json:"oauth" yaml:"oauth" mapstructure:"oauth"`
	Log        types.LogConfig        `json:"log" yaml:"log" mapstructure:"log"`
	Middleware types.MiddlewareConfig `json:"middleware" yaml:"middleware" mapstructure:"middleware"`
	Casdoor    types.CasdoorConfig    `json:"casdoor" yaml:"casdoor" mapstructure:"casdoor"`
	AliyunOSS  types.AliyunOSSConfig  `json:"aliyun_oss" yaml:"aliyun_oss" mapstructure:"aliyun_oss"`
}
