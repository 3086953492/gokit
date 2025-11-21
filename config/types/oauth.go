package types

import "time"

// OAuthConfig OAuth 认证配置
type OAuthConfig struct {
	Secret             string        `json:"secret" yaml:"secret" mapstructure:"secret"`
	AuthCodeExpire     time.Duration `json:"auth_code_expire" yaml:"auth_code_expire" mapstructure:"auth_code_expire"`
	AccessTokenExpire  time.Duration `json:"access_token_expire" yaml:"access_token_expire" mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `json:"refresh_token_expire" yaml:"refresh_token_expire" mapstructure:"refresh_token_expire"`
}
