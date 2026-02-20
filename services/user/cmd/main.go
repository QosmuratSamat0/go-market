package main

import (
	"log/slog"
	"os"

	"github.com/go-market/pkg/logger/handlers/slogpretty"
	"github.com/go-market/services/user/internal/app"
	"github.com/go-market/services/user/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	log := setupPrettyLogger()

	log.Info("info level starting up", slog.String("env", cfg.Env))
	log.Debug("debug level starting up")

	application, err := app.NewApp(*cfg, log)
	if err != nil {
		log.Error("failed to initialize app", slog.Any("err", err))
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		log.Error("app stopped with error", slog.Any("err", err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
