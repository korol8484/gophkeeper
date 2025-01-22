package serve

import (
	"database/sql"
	"fmt"
	"github.com/korol8484/gophkeeper/internal/server/api/secret"
	"github.com/korol8484/gophkeeper/internal/server/api/user"
	"github.com/korol8484/gophkeeper/internal/server/app"
	"github.com/korol8484/gophkeeper/internal/server/db"
	"github.com/korol8484/gophkeeper/internal/server/logger"
	"github.com/korol8484/gophkeeper/internal/server/secret/add"
	"github.com/korol8484/gophkeeper/internal/server/secret/get"
	"github.com/korol8484/gophkeeper/internal/server/secret/list"
	sercretrepository "github.com/korol8484/gophkeeper/internal/server/secret/repository"
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

			repo := sercretrepository.NewSecretRepository(repoDb)

			serviceAdd := add.NewSecretService(repo)
			addHandler := secret.NewSecretAddHandler(serviceAdd)

			serviceList := list.NewSecretServiceList(repo)
			listHandler := secret.NewListHandler(serviceList)

			serviceGet := get.NewSecretServiceList(repo)
			getHandler := secret.NewGetHandler(serviceGet)

			httpSvc.AddHandler(addHandler.RegisterRoutes(jwtSvc))
			httpSvc.AddHandler(listHandler.RegisterRoutes(jwtSvc))
			httpSvc.AddHandler(getHandler.RegisterRoutes(jwtSvc))

			errCh := make(chan error, 1)
			oss, stop := make(chan os.Signal, 1), make(chan struct{}, 1)
			signal.Notify(oss, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

			go func() {
				if err = httpSvc.Run(*debug); err != nil {
					errCh <- err
				}
			}()

			for {
				select {
				case <-oss:
					httpSvc.Stop()
					return nil
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
