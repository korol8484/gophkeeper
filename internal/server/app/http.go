package app

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/korol8484/gophkeeper/internal/server/app/middlewares"
	"github.com/korol8484/gophkeeper/pkg"
	"go.uber.org/zap"
	"net/http"
)

type middlewareHandler func(http.Handler) http.Handler
type registerHandler func(mux *chi.Mux)

type App struct {
	httpServer *http.Server
	cfg        *Config
	log        *zap.Logger

	middlewares []middlewareHandler
	handlers    []registerHandler
}

// NewApp - factory
func NewApp(
	cfg *Config,
	log *zap.Logger,
) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

// Run - start http server
func (a *App) Run(debug bool) error {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		middlewares.NewLogging(a.log, 2).LoggingMiddleware,
	)

	for _, m := range a.middlewares {
		router.Use(m)
	}

	a.httpServer = &http.Server{
		Addr:    a.cfg.Listen,
		Handler: router,
	}

	for _, hr := range a.handlers {
		hr(router)
	}

	errCh := make(chan error, 1)

	go func() {
		err := a.httpServer.ListenAndServeTLS(a.cfg.Pem, a.cfg.Key)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		close(errCh)
	}()

	return <-errCh
}

// AddMiddlewares -
func (a *App) AddMiddlewares(handler middlewareHandler) {
	a.middlewares = append(a.middlewares, handler)
}

// AddHandler -
func (a *App) AddHandler(handler registerHandler) {
	a.handlers = append(a.handlers, handler)
}

// Stop -
func (a *App) Stop() {
	ctx, shutdown := context.WithTimeout(context.Background(), pkg.TimeOut)
	defer shutdown()

	_ = a.httpServer.Shutdown(ctx)
}
