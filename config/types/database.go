package types

type DatabaseConfig struct {
	Host      string `json:"host" yaml:"host" mapstructure:"host"`
	Port      int    `json:"port" yaml:"port" mapstructure:"port"`
	User      string `json:"user" yaml:"user" mapstructure:"user"`
	Password  string `json:"password" yaml:"password" mapstructure:"password"`
	DBName    string `json:"dbname" yaml:"dbname" mapstructure:"dbname"`
	Charset   string `json:"charset" yaml:"charset" mapstructure:"charset"`
	ParseTime bool   `json:"parseTime" yaml:"parseTime" mapstructure:"parseTime"`
	Loc       string `json:"loc" yaml:"loc" mapstructure:"loc"`
}
