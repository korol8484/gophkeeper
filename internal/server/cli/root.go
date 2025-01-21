package cli

import (
	"errors"
	"github.com/korol8484/gophkeeper/internal/server/cli/migrate"
	"github.com/korol8484/gophkeeper/internal/server/cli/serve"
	"github.com/korol8484/gophkeeper/internal/server/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func NewRootCommand() *cobra.Command {
	cfgFile := util.ToPtr("")
	debug := util.ToPtr(false)

	cmd := &cobra.Command{
		Use:           "app",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "App cli",
		PersistentPreRunE: func(*cobra.Command, []string) error {
			if cfgFile == nil || *cfgFile == "" {
				return errors.New("no configuration file provided")
			}

			if _, err := os.Stat(*cfgFile); err != nil {
				return err
			}

			if absPath, err := filepath.Abs(*cfgFile); err == nil {
				*cfgFile = absPath

				if err = os.Chdir(filepath.Dir(absPath)); err != nil {
					return err
				}
			}

			viper.SetConfigFile(*cfgFile)
			viper.SetEnvPrefix("SERVER")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()

			if err := viper.ReadInConfig(); err == nil {
				return err
			}

			return nil
		},
	}

	f := cmd.PersistentFlags()

	f.StringVarP(cfgFile, "config", "c", ".app.yaml", "config file (default is .app.yaml)")
	f.BoolVarP(debug, "debug", "d", false, "debug mode")

	cmd.AddCommand(serve.NewServeCommand(debug))
	cmd.AddCommand(migrate.NewMigrateCommand())

	return cmd
}
