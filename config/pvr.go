package config

type Pvr struct {
	Enabled bool
	Type    string
	URL     string
	ApiKey  string `mapstructure:"api_key"`
}
