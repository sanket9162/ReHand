package main

import (
	"log"
	"net/http"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/handler"
)

type application struct {
	handler *handler.Handler
	config  *config.Config
}

func main() {

	cfg := config.LoadConfig()
	h := handler.NewHandler()

	app := application{
		handler: h,
		config:  cfg,
	}
	err := http.ListenAndServe(":"+app.config.Port, app.router())
	if err != nil {
		log.Fatalf("server failed %v", err)
	}
}
