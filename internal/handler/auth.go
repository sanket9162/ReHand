package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sanket9162/rehand/internal/middleware"
	"github.com/sanket9162/rehand/internal/models"
	"github.com/sanket9162/rehand/internal/repository"
	"github.com/sanket9162/rehand/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	var input SignUpRequest
	if err := utils.ReadJSON(w, r, &input); err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeMalformedJSON, err.Error())
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("failed to generate password hash", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	created, err := h.Repo.CreateUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			_ = utils.ErrorJSON(w, http.StatusConflict, utils.CodeConflict, "email already exists", "email")
			return
		}
		h.logger.Error("failed to create user", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}
	if err := utils.WriteJSON(w, http.StatusCreated, created); err != nil {
		h.logger.Error("failed to write signup response", "request_id", reqID, "error", err)
	}

}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	var input SignInRequest
	if err := utils.ReadJSON(w, r, &input); err != nil {
		_ = utils.ErrorJSON(w, http.StatusBadRequest, utils.CodeMalformedJSON, err.Error())
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

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

	user, err := h.Repo.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "invalid email or password")
			return
		}
		h.logger.Error("failed to get user by email", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "invalid email or password")
			return
		}
		h.logger.Error("failed to compare password hash", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	tokens, err := h.jwt.GenerateTokenPair(&utils.JwtUser{ID: user.ID})
	if err != nil {
		h.logger.Error("failed to generate tokens", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	http.SetCookie(w, h.jwt.GetAccessTokenCookie(tokens.Token))
	http.SetCookie(w, h.jwt.GetRefreshCookie(tokens.RefreshToken))

	resp := map[string]any{
		"access_token":  tokens.Token,
		"refresh_token": tokens.RefreshToken,
		"user":          user,
	}
	if err := utils.WriteJSON(w, http.StatusOK, resp); err != nil {
		h.logger.Error("failed to write signin response", "request_id", reqID, "error", err)
	}
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "missing refresh token")
		return
	}

	claims, err := h.jwt.ValidateToken(cookie.Value, utils.TokenTypeRefresh)
	if err != nil {
		_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "invalid or expired refresh token")
		return
	}

	user, err := h.Repo.GetUserByID(r.Context(), claims.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "user not found")
			return
		}
		h.logger.Error("failed to get user during refresh", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	tokens, err := h.jwt.GenerateTokenPair(&utils.JwtUser{ID: user.ID})
	if err != nil {
		h.logger.Error("failed to generate refresh tokens", "request_id", reqID, "error", err)
		_ = utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error")
		return
	}

	http.SetCookie(w, h.jwt.GetAccessTokenCookie(tokens.Token))
	http.SetCookie(w, h.jwt.GetRefreshCookie(tokens.RefreshToken))

	if err := utils.WriteJSON(w, http.StatusOK, tokens); err != nil {
		h.logger.Error("failed to write refresh response", "request_id", reqID, "error", err)
	}
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	http.SetCookie(w, h.jwt.GetExpiredAccessTokenCookie())
	http.SetCookie(w, h.jwt.GetExpiredRefreshCookie())

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "signed out"}); err != nil {
		h.logger.Error("failed to write signout response", "request_id", reqID, "error", err)
	}
}
