package types

type RedisConfig struct {
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Port     int    `json:"port" yaml:"port" mapstructure:"port"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	DB       int    `json:"db" yaml:"db" mapstructure:"db"`
}
