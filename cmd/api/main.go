package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	var app application
	err := http.ListenAndServe(":8000", app.router())
	if err != nil {
		log.Fatalf("server failed %v", err)
	}
}
