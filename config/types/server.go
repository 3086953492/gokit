package types

type ServerConfig struct {
	BaseURL string `json:"base_url" yaml:"base_url" mapstructure:"base_url"`
	Port    int    `json:"port" yaml:"port" mapstructure:"port"`
	Mode    string `json:"mode" yaml:"mode" mapstructure:"mode"`
	ID      int    `json:"id" yaml:"id" mapstructure:"id"`
}
