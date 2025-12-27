package types

type ServerConfig struct {
	BaseURL     string `json:"base_url" yaml:"base_url" mapstructure:"base_url"`
	Port        int    `json:"port" yaml:"port" mapstructure:"port"`
	FrontendURL string `json:"frontend_url" yaml:"frontend_url" mapstructure:"frontend_url"`
	Mode        string `json:"mode" yaml:"mode" mapstructure:"mode"`
	ID          int    `json:"id" yaml:"id" mapstructure:"id"`
	Domain      string `json:"domain" yaml:"domain" mapstructure:"domain"`
}
