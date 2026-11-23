package main

import (
	"log"
	"net/http"

	"github.com/sanket9162/rehand/internal/handler"
)

type application struct {
	handler *handler.Handler
}

func main() {
	var app application
	err := http.ListenAndServe(":8000", app.router())
	if err != nil {
		log.Fatalf("server failed %v", err)
	}
}
