package app

import (
	"time"
)

type Config struct {
	Listen         string        `mapstructure:"listen"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
	Key            string        `mapstructure:"key"`
	Pem            string        `mapstructure:"pem"`
}
