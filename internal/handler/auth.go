package handler

import (
	"net/http"

	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/utils"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	if err := utils.WriteJSON(w, http.StatusOK, reqID); err != nil {
		h.logger.Error("failed to write create listing response", "request_id", reqID, "error", err)
	}
}
