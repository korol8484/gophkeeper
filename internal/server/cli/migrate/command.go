package migrate

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/korol8484/gophkeeper/internal/server/cli/util"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigrateCommand() *cobra.Command {
	up := util.ToPtr(false)
	down := util.ToPtr(false)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run service migrations",
		RunE: func(*cobra.Command, []string) error {
			cfg, err := newConfig()
			if err != nil {
				return err
			}

			if *up && *down {
				return errors.New("use only one, up or down")
			}

			m, err := migrate.New(
				"file://internal/server/migrations",
				cfg.Db.Dsn,
			)
			if err != nil {
				return err
			}

			if *up {
				if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return err
				}

				fmt.Println("migrate up success")
			}

			if *down {
				if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return err
				}

				fmt.Println("migrate down success")
			}

			return nil
		},
	}

	f := cmd.PersistentFlags()
	f.BoolVarP(up, "up", "u", false, "migrate up")
	f.BoolVarP(down, "down", "n", false, "migrate down")

	return cmd
}
