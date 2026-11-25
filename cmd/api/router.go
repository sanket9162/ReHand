package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) router() http.Handler {
	router := chi.NewRouter()

	router.Use(app.middleware.RequestID)
	router.Use(app.middleware.Recoverer)

	router.Get("/healthz", app.handler.HealthCheck)
	router.Get("/listings", app.handler.GetListings)
	router.Delete("/listings/{id}", app.handler.DeleteListing)

	return router
}
