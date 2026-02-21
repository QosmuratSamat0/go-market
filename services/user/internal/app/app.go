package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-market/services/user/internal/config"
	userHTTP "github.com/go-market/services/user/internal/derivery/http"
	"github.com/go-market/services/user/internal/repository/postgres"
	"github.com/go-market/services/user/internal/service"
)

type App struct {
	server *http.Server
	log    *slog.Logger
}

func NewApp(cfg config.Config, log *slog.Logger) (*App, error) {
	const op = "app.NewApp"

	logger := log.With(slog.String("op", op))

	repo, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to init postgres", slog.String("op", op), slog.Any("err", err))
		return nil, err
	}
	svc := service.New(repo)

	userHandler := userHTTP.New(log, svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		userHTTP.RegisterUserRoutes(r, userHandler, cfg.SecretKey)
	})

	server := &http.Server{
		Addr:         cfg.HTTPAddr.Address,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{server: server, log: log}, nil
}

func (a *App) Run() error {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("listen failed", slog.Any("err", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		panic(err)
	}

	return nil
}
