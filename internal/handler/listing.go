package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/repository"
	"github.com/sanket9162/rehand/internal/utils"
)

func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.Repo.Listings(r.Context())
	if err != nil {
		h.logger.Error("failed to get listings", "error", err)
		_ = utils.ErrorJSON(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, listings); err != nil {
		h.logger.Error("failed to write listings response", "error", err)
	}
}

func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = utils.ErrorJSON(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}

	newId, err := uuid.Parse(id)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteListing(r.Context(), newId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, errors.New("listing not found"), http.StatusNotFound)
			return
		}
		h.logger.Error("failed to delete listing", "error", err)
		_ = utils.ErrorJSON(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"deleted_id": id}); err != nil {
		h.logger.Error("failed to write delete response", "error", err)
	}
}
