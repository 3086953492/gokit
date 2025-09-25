package types

type CasdoorConfig struct {
	Host   string `json:"host"`
	Client Client `json:"client"`
}

type Client struct {
	Endpoint         string `json:"endpoint"`
	ClientId         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	Certificate      string `json:"certificate"`
	OrganizationName string `json:"organization_name"`
	ApplicationName  string `json:"application_name"`
}
