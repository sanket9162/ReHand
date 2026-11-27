package middleware

import (
	"log/slog"

	"github.com/sanket9162/rehand/internal/utils"
)

type Middleware struct {
	Logger *slog.Logger
	JWT    *utils.JWT
}

func NewMiddleware(logger *slog.Logger, jwt *utils.JWT) *Middleware {
	return &Middleware{
		Logger: logger,
		JWT:    jwt,
	}
}
