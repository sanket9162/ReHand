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
	router.NotFound(app.handler.NotFound)
	router.MethodNotAllowed(app.handler.MethodNotAllowed)

	// Public routes
	router.Get("/listings", app.handler.GetListings)
	router.Get("/listings/{id}", app.handler.GetListingByID)
	router.Post("/signup", app.handler.SignUp)
	router.Post("/signin", app.handler.SignIn)
	router.Post("/refresh", app.handler.RefreshToken)
	router.Post("/signout", app.handler.SignOut)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(app.middleware.RequireAuth)

		r.Post("/listings", app.handler.CreateListing)
		r.Patch("/listings/{id}", app.handler.UpdateListing)
		r.Delete("/listings/{id}", app.handler.DeleteListing)
	})

	return router
}
