package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sanket9162/rehand/internal/utils"
)

func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.Repo.Listings(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, listings)
}

func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, id)
}
