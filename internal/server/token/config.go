package token

import "time"

type Config struct {
	Secret string        `mapstructure:"secret"`
	Name   string        `mapstructure:"name"`
	Expire time.Duration `mapstructure:"expire"`
}
