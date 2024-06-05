package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/sahidrahman404/todos/internal/database"
	"github.com/sahidrahman404/todos/internal/env"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	dsn      string
	httpPort int
}

type application struct {
	db     *database.DB
	logger *slog.Logger
	config config
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.dsn = env.GetString("DB_DSN", "postgres://todo:password@172.17.0.2:5432/todo?sslmode=disable")

	db, err := database.New(cfg.dsn)
	if err != nil {
		return err
	}

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
	}

	return app.serveHTTP()
}
