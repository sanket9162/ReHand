package main

import (
	"log"
	"net/http"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/database"
	"github.com/sanket9162/rehand/internal/handler"
	"github.com/sanket9162/rehand/internal/repository/dbrepo"
)

type application struct {
	handler *handler.Handler
	config  *config.Config
}

func main() {
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

	h := handler.NewHandler(repo)

	app := application{
		handler: h,
		config:  cfg,
	}
	err = http.ListenAndServe(":"+app.config.Port, app.router())
	if err != nil {
		log.Fatalf("server failed %v", err)
	}
}
