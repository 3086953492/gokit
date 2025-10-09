package types

import "time"

type JWTConfig struct {
	Secret string        `json:"secret" yaml:"secret" mapstructure:"secret"`
	Expire time.Duration `json:"expire" yaml:"expire" mapstructure:"expire"`
}
