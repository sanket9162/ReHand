package handler

import (
	"net/http"

	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/utils"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	if err := utils.WriteJSON(w, http.StatusOK, reqID); err != nil {
		h.logger.Error("failed to write create listing response", "request_id", reqID, "error", err)
	}
}
