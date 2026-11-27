package handler

import (
	"log/slog"
	"net/http"

	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/repository"
	"github.com/sanket9162/rehand/internal/utils"
)

type Handler struct {
	Repo   repository.DatabaseRepo
	logger *slog.Logger
	jwt    *utils.JWT
}

func NewHandler(repo repository.DatabaseRepo, logger *slog.Logger, jwt *utils.JWT) *Handler {
	return &Handler{
		Repo:   repo,
		logger: logger,
		jwt:    jwt,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.Repo.Ping(r.Context()); err != nil {
		h.logger.Error("healthcheck database ping failed",
			"request_id", middleware.GetRequestID(r.Context()),
			"error", err,
		)
		_ = utils.ErrorJSON(w, http.StatusServiceUnavailable, utils.CodeInternalError, "database unavailable")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		h.logger.Error("failed to write healthcheck response",
			"request_id", middleware.GetRequestID(r.Context()),
			"error", err,
		)
	}
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	if err := utils.ErrorJSON(w, http.StatusNotFound, utils.CodeNotFound, "route not found"); err != nil {
		h.logger.Error("failed to write not found response",
			"request_id", middleware.GetRequestID(r.Context()),
			"error", err,
		)
	}
}

func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if err := utils.ErrorJSON(w, http.StatusMethodNotAllowed, utils.CodeMethodNotAllowed, "method not allowed"); err != nil {
		h.logger.Error("failed to write method not allowed response",
			"request_id", middleware.GetRequestID(r.Context()),
			"error", err,
		)
	}
}
