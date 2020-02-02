package config

import "time"

type Pvr struct {
	Enabled      bool
	Type         string
	URL          string
	ApiKey       string        `mapstructure:"api_key"`
	RetryDaysAge time.Duration `mapstructure:"retry_days_age"`
}
