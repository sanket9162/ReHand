package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sanket9162/rehand/internal/middleware"
)

func (app *application) router() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/healthz", app.handler.HealthCheck)
	router.Get("/listings", app.handler.GetListings)
	router.Delete("/listings/{id}", app.handler.DeleteListing)

	return router
}
