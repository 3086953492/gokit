package types

import "time"

// OAuthConfig OAuth 认证配置
type OAuthConfig struct {
	Secret             string        `json:"secret" yaml:"secret" mapstructure:"secret"`
	AuthCodeExpire     time.Duration `json:"auth_code_expire" yaml:"auth_code_expire" mapstructure:"auth_code_expire"`
	AccessTokenExpire  time.Duration `json:"access_token_expire" yaml:"access_token_expire" mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `json:"refresh_token_expire" yaml:"refresh_token_expire" mapstructure:"refresh_token_expire"`

	// Subject 相关配置（用于生成 OIDC sub）
	Subject SubjectConfig `json:"subject" yaml:"subject" mapstructure:"subject"`
}

// SubjectConfig OIDC Subject 标识符配置
type SubjectConfig struct {
	// Secret HMAC 密钥，用于生成稳定的 sub，必须至少 32 字节且长期稳定
	Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`

	// Length sub 的长度（不含前缀），0 表示不截断（默认 43 字符）
	Length int `json:"length" yaml:"length" mapstructure:"length"`

	// Prefix sub 的前缀，例如 "u_"，便于区分来源
	Prefix string `json:"prefix" yaml:"prefix" mapstructure:"prefix"`
}
