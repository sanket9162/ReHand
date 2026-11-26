package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/models"
	"github.com/sanket9162/rehand/internal/repository"
	"github.com/sanket9162/rehand/internal/utils"
)

type createListingRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
	City        string `json:"city" validate:"required"`
}

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

	var input createListingRequest
	if err := utils.ReadJSON(w, r, &input); err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeMalformedJSON, err.Error())
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.City = strings.TrimSpace(input.City)

	if err := utils.ValidateStruct(input); err != nil {
		var valErr *utils.ValidationError
		if errors.As(err, &valErr) {
			_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeValidationFailed, valErr.Message, valErr.Field)
			return
		}
		h.logger.Error("unexpected validation error", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	listing := &models.Listing{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		City:        input.City,
	}

	created, err := h.Repo.CreateListing(r.Context(), listing)
	if err != nil {
		h.logger.Error("failed to create listing", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}
	headers := http.Header{"Location": {"/listings/" + created.ID.String()}}

	if err := utils.WriteJSON(w, http.StatusCreated, created, headers); err != nil {
		h.logger.Error("failed to write create listing response", "request_id", reqID, "error", err)
	}
}
