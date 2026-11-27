package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/utils"
)

const UserIDKey contextKey = "user_id"

func (mw *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := mw.JWT.GetTokenFromHeader(r)
		if err != nil {
			cookie, cookieErr := r.Cookie("access_token")
			if cookieErr != nil || cookie.Value == "" {
				_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "missing or invalid authorization header")
				return
			}
			tokenStr = cookie.Value
		}

		claims, err := mw.JWT.ValidateToken(tokenStr, utils.TokenTypeAccess)
		if err != nil {
			_ = utils.ErrorJSON(w, http.StatusUnauthorized, utils.CodeUnauthenticated, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
