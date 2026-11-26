package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/repository"
	"github.com/sanket9162/rehand/internal/utils"
)

func (h *Handler) GetListings(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	listings, err := h.Repo.Listings(r.Context())
	if err != nil {
		h.logger.Error("failed to get listings", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, listings); err != nil {
		h.logger.Error("failed to write listings response", "request_id", reqID, "error", err)
	}
}

func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	id := chi.URLParam(r, "id")
	if id == "" {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeInvalidID, "id is required", "id")
		return
	}

	newId, err := uuid.Parse(id)
	if err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeInvalidID, "invalid id", "id")
		return
	}

	err = h.Repo.DeleteListing(r.Context(), newId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, http.StatusNotFound, utils.CodeNotFound, "listing not found")
			return
		}
		h.logger.Error("failed to delete listing", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"deleted_id": id}); err != nil {
		h.logger.Error("failed to write delete response", "request_id", reqID, "error", err)
	}
}

func (h *Handler) CreateListings(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "ok"}); err != nil {
		h.logger.Error("failed to write delete response", "request_id", reqID, "error", err)
	}
}
