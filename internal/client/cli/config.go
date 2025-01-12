package cli

import (
	"github.com/korol8484/gophkeeper/internal/client/service"
	"github.com/spf13/viper"
	"strings"
)

type config struct {
	client *service.Config `mapstructure:"client"`
}

func newConfig() (*config, error) {
	cfg := &config{
		client: &service.Config{
			ServiceHost: "https://localhost:8091",
		},
	}

	viper.SetEnvPrefix("CLIENT")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		return nil, err
	}

	return cfg, nil
}
