package main

import (
	"log"
	"net/http"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/database"
	"github.com/sanket9162/rehand/internal/handler"
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

	h := handler.NewHandler()

	app := application{
		handler: h,
		config:  cfg,
	}
	err = http.ListenAndServe(":"+app.config.Port, app.router())
	if err != nil {
		log.Fatalf("server failed %v", err)
	}
}
