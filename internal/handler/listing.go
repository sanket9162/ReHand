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

type updateListingRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1"`
	Description *string `json:"description" validate:"omitempty,min=1"`
	Price       *int64  `json:"price" validate:"omitempty,gt=0"`
	City        *string `json:"city" validate:"omitempty,min=1"`
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

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "unauthorized")
		return
	}

	newId, err := uuid.Parse(id)
	if err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeInvalidID, "invalid id", "id")
		return
	}

	err = h.Repo.DeleteListing(r.Context(), newId, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, http.StatusNotFound, utils.CodeNotFound, "listing not found")
			return
		}
		if errors.Is(err, repository.ErrForbidden) {
			_ = utils.ErrorJSON(w, http.StatusForbidden, utils.CodeForbidden, "you do not own this listing")
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

func (h *Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "unauthorized")
		return
	}

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
    UserID:      userID,
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

func (h *Handler) GetListingByID(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeInvalidID, "invalid id", "id")
		return
	}
	listing, err := h.Repo.GetListingByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, http.StatusNotFound, utils.CodeNotFound, "listing not found")
			return
		}
		h.logger.Error("failed to get listing", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}
	if err := utils.WriteJSON(w, http.StatusOK, listing); err != nil {
		h.logger.Error("failed to write listing response", "request_id", reqID, "error", err)
	}
}
func (h *Handler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeInvalidID, "invalid id", "id")
		return
	}
	var input updateListingRequest
	if err := utils.ReadJSON(w, r, &input); err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeMalformedJSON, err.Error())
		return
	}
	if input.Title != nil {
		trimmed := strings.TrimSpace(*input.Title)
		input.Title = &trimmed
	}
	if input.Description != nil {
		trimmed := strings.TrimSpace(*input.Description)
		input.Description = &trimmed
	}
	if input.City != nil {
		trimmed := strings.TrimSpace(*input.City)
		input.City = &trimmed
	}
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
	existing, err := h.Repo.GetListingByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, http.StatusNotFound, utils.CodeNotFound, "listing not found")
			return
		}
		h.logger.Error("failed to load listing for update", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}
	if existing.UserID != userID {
		_ = utils.ErrorJSON(w, http.StatusForbidden, utils.CodeForbidden, "you do not own this listing")
		return
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Price != nil {
		existing.Price = *input.Price
	}
	if input.City != nil {
		existing.City = *input.City
	}
	updated, err := h.Repo.UpdateListing(r.Context(), existing)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			_ = utils.ErrorJSON(w, http.StatusForbidden, utils.CodeForbidden, "you do not own this listing")
			return
		}
		h.logger.Error("failed to update listing", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}
	if err := utils.WriteJSON(w, http.StatusOK, updated); err != nil {
		h.logger.Error("failed to write update listing response", "request_id", reqID, "error", err)
	}
}