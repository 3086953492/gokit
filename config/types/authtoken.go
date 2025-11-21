package types

import "time"

// AuthTokenConfig 认证令牌配置
type AuthTokenConfig struct {
	Secret        string        `json:"secret" yaml:"secret" mapstructure:"secret"`
	AccessExpire  time.Duration `json:"access_expire" yaml:"access_expire" mapstructure:"access_expire"`
	RefreshExpire time.Duration `json:"refresh_expire" yaml:"refresh_expire" mapstructure:"refresh_expire"`
	Issuer        string        `json:"issuer" yaml:"issuer" mapstructure:"issuer"`
}
