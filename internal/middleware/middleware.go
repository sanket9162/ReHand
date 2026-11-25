package middleware

import "log/slog"

type Middleware struct {
	Logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{
		Logger: logger,
	}
}
