package types

type CasdoorConfig struct {
	Host   string `json:"host" yaml:"host" mapstructure:"host"`
	Client Client `json:"client" yaml:"client" mapstructure:"client"`
}

type Client struct {
	Endpoint         string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	ClientId         string `json:"client_id" yaml:"client_id" mapstructure:"client_id"`
	ClientSecret     string `json:"client_secret" yaml:"client_secret" mapstructure:"client_secret"`
	Certificate      string `json:"certificate" yaml:"certificate" mapstructure:"certificate"`
	OrganizationName string `json:"organization_name" yaml:"organization_name" mapstructure:"organization_name"`
	ApplicationName  string `json:"application_name" yaml:"application_name" mapstructure:"application_name"`
}
