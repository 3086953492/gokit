package types

import "time"

// AuthTokenConfig 认证令牌配置
type AuthTokenConfig struct {
	AccessTokenSecret  string        `json:"access_token_secret" yaml:"access_token_secret" mapstructure:"access_token_secret"`
	AccessTokenExpire  time.Duration `json:"access_token_expire" yaml:"access_token_expire" mapstructure:"access_token_expire"`
	RefreshTokenSecret string        `json:"refresh_token_secret" yaml:"refresh_token_secret" mapstructure:"refresh_token_secret"`
	RefreshTokenExpire time.Duration `json:"refresh_token_expire" yaml:"refresh_token_expire" mapstructure:"refresh_token_expire"`
	Issuer             string        `json:"issuer" yaml:"issuer" mapstructure:"issuer"`
}
