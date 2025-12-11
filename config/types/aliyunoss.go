package types

type AliyunOSSConfig struct {
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Endpoint        string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Bucket          string `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	Domain          string `json:"domain" yaml:"domain" mapstructure:"domain"` // 可选：自定义域名
}
