package config

import "time"

type Pvr struct {
	Type         string
	URL          string
	ApiKey       string       `mapstructure:"api_key"`
	RetryDaysAge RetryDaysAge `mapstructure:"retry_days_age"`
}

type RetryDaysAge struct {
	Missing time.Duration
	Cutoff  time.Duration
}
