package types

type GoauthConfig struct {
	FrontendBaseURL    string `json:"frontend_base_url" yaml:"frontend_base_url" mapstructure:"frontend_base_url"`
	BackendBaseURL     string `json:"backend_base_url" yaml:"backend_base_url" mapstructure:"backend_base_url"`
	ClientID           string `json:"client_id" yaml:"client_id" mapstructure:"client_id"`
	ClientSecret       string `json:"client_secret" yaml:"client_secret" mapstructure:"client_secret"`
	RedirectURI        string `json:"redirect_uri" yaml:"redirect_uri" mapstructure:"redirect_uri"`
	AccessTokenSecret  string `json:"access_token_secret" yaml:"access_token_secret" mapstructure:"access_token_secret"`
	RefreshTokenSecret string `json:"refresh_token_secret" yaml:"refresh_token_secret" mapstructure:"refresh_token_secret"`
}
