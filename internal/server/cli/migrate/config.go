package migrate

import (
	"fmt"
	"github.com/korol8484/gophkeeper/internal/server/db"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Db *db.Config `mapstructure:"db"`
}

func newConfig() (*Config, error) {
	cfg := &Config{
		Db: &db.Config{
			MaxIdleConn:     1,
			MaxOpenConn:     10,
			MaxLifetimeConn: time.Minute * 3,
		},
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %w", err)
	}

	return cfg, nil
}
