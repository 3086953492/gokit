package config

import "github.com/3086953492/YaBase/config/types"

type Config struct {
	Server     types.ServerConfig     `mapstructure:"server"`
	Database   types.DatabaseConfig   `mapstructure:"database"`
	Redis      types.RedisConfig      `mapstructure:"redis"`
	JWT        types.JWTConfig        `mapstructure:"jwt"`
	Log        types.LogConfig        `mapstructure:"log"`
	Middleware types.MiddlewareConfig `mapstructure:"middleware"`
	Casdoor    types.CasdoorConfig    `mapstructure:"casdoor"`
}
