package handler

import (
	"fmt"
	"net/http"

	"github.com/sanket9162/rehand/internal/utils"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		fmt.Println("error writing json:", err)
	}
}
