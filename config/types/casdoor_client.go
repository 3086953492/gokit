package types

type CasdoorClientConfig struct {
	Endpoint         string `json:"endpoint"`
	ClientId         string `json:"clientId"`
	ClientSecret     string `json:"clientSecret"`
	Certificate      string `json:"certificate"`
	OrganizationName string `json:"organizationName"`
	ApplicationName  string `json:"applicationName"`
}
