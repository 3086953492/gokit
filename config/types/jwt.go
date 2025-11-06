package types

import "time"

type JWTConfig struct {
	Secret        string        `json:"secret" yaml:"secret" mapstructure:"secret"`
	Expire        time.Duration `json:"expire" yaml:"expire" mapstructure:"expire"`
	RefreshExpire time.Duration `json:"refresh_expire" yaml:"refresh_expire" mapstructure:"refresh_expire"`
	Issuer        string        `json:"issuer" yaml:"issuer" mapstructure:"issuer"`
}
