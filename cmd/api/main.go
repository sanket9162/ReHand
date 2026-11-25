package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/database"
	"github.com/sanket9162/rehand/internal/handler"
	"github.com/sanket9162/rehand/internal/repository/dbrepo"
)

type application struct {
	logger  *slog.Logger
	handler *handler.Handler
	config  *config.Config
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
		log.Fatalf("database connection failed %v", err)
	}
	defer db.SQL.Close()

	// Initialize application
	repo := dbrepo.NewPostgresDBRepo(db.SQL, cfg)
	h := handler.NewHandler(repo, logger)

	app := application{
		handler: h,
		config:  cfg,
	}
	if err = app.server(); err != nil {
		log.Fatalf("server failed %v", err)
	}
}
