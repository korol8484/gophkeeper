package serve

import (
	"database/sql"
	"fmt"
	"github.com/korol8484/gophkeeper/internal/server/api/user"
	"github.com/korol8484/gophkeeper/internal/server/app"
	"github.com/korol8484/gophkeeper/internal/server/db"
	"github.com/korol8484/gophkeeper/internal/server/logger"
	"github.com/korol8484/gophkeeper/internal/server/token"
	"github.com/korol8484/gophkeeper/internal/server/user/auth"
	"github.com/korol8484/gophkeeper/internal/server/user/repository"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func NewServeCommand(debug *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run as service",
		RunE: func(*cobra.Command, []string) error {
			cfg, err := newConfig()
			if err != nil {
				return err
			}

			logSvc, err := logger.NewLogger(*debug)
			if err != nil {
				return fmt.Errorf("can't create loging service: %w", err)
			}

			repoDb, err := db.NewPgDB(cfg.Db)
			if err != nil {
				return fmt.Errorf("can't create repository db connection: %w", err)
			}

			defer func(repoDb *sql.DB) {
				_ = repoDb.Close()
			}(repoDb)

			jwtSvc := token.NewJwtService(cfg.Token)
			userDb := repository.NewDBStore(repoDb)
			userSvc := auth.NewService(userDb)
			userHandler := user.NewAuthHandler(userSvc, jwtSvc)

			httpSvc := app.NewApp(cfg.Http, logSvc)
			httpSvc.AddHandler(userHandler.RegisterRoutes())

			errCh := make(chan error, 1)
			oss, stop := make(chan os.Signal, 1), make(chan struct{}, 1)
			signal.Notify(oss, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

			go func() {
				<-oss

				stop <- struct{}{}
			}()

			go func() {
				if err = httpSvc.Run(*debug); err != nil {
					errCh <- err
				}
			}()

			for {
				select {
				case e := <-errCh:
					return e
				case <-stop:
					httpSvc.Stop()
					return nil
				}
			}
		},
	}
}
