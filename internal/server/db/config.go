package db

import (
	"time"
)

type Config struct {
	Dsn             string        `mapstructure:"dsn"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
	MaxLifetimeConn time.Duration `mapstructure:"max_lifetime_conn"`
}
