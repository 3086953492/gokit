package config

import "github.com/3086953492/gokit/config/types"

type Config struct {
	Server     types.ServerConfig     `json:"server" yaml:"server" mapstructure:"server"`
	Database   types.DatabaseConfig   `json:"database" yaml:"database" mapstructure:"database"`
	Redis      types.RedisConfig      `json:"redis" yaml:"redis" mapstructure:"redis"`
	AuthToken  types.AuthTokenConfig  `json:"auth_token" yaml:"auth_token" mapstructure:"auth_token"`
	Log        types.LogConfig        `json:"log" yaml:"log" mapstructure:"log"`
	Middleware types.MiddlewareConfig `json:"middleware" yaml:"middleware" mapstructure:"middleware"`
	Casdoor    types.CasdoorConfig    `json:"casdoor" yaml:"casdoor" mapstructure:"casdoor"`
}
