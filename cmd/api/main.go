package main

import (
	"log/slog"
	"os"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/database"
	"github.com/sanket9162/rehand/internal/handler"
	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/repository/dbrepo"
)

type application struct {
	logger     *slog.Logger
	handler    *handler.Handler
	middleware *middleware.Middleware
	config     *config.Config
}

func main() {
	// initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Load config
	cfg := config.LoadConfig()

	// initialize database connection
	db, err := database.ConnectPostgres(cfg.DSN)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.SQL.Close()

	// Initialize application
	repo := dbrepo.NewPostgresDBRepo(db.SQL)
	h := handler.NewHandler(repo, logger)
	mw := middleware.NewMiddleware(logger)

	app := application{
		handler:    h,
		middleware: mw,
		config:     cfg,
		logger:     logger,
	}
	if err = app.server(); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
