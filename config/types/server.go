package types

type ServerConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	ID      int    `mapstructure:"id"`
}
