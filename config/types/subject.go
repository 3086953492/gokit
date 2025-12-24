package types

type SubjectConfig struct {
	Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`
	Length int    `json:"length" yaml:"length" mapstructure:"length"`
	Prefix string `json:"prefix" yaml:"prefix" mapstructure:"prefix"`
}
