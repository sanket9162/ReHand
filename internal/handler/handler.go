package handler

import (
	"log/slog"
	"net/http"

	"github.com/sanket9162/rehand/internal/repository"
	"github.com/sanket9162/rehand/internal/utils"
)

type Handler struct {
	Repo   repository.DatabaseRepo
	logger *slog.Logger
}

func NewHandler(repo repository.DatabaseRepo, logger *slog.Logger) *Handler {
	return &Handler{
		Repo:   repo,
		logger: logger,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
