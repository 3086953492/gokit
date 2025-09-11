package types

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	ID   int    `mapstructure:"id"`
}
