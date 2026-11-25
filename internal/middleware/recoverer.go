// internal/middleware/recoverer.go
package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/sanket9162/rehand/internal/utils"
)

func (mw *Middleware) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				// http.ErrAbortHandler is a special sentinel panic in net/http;
				// it should not be recovered or logged as an error.
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				reqID := GetRequestID(r.Context())
				mw.Logger.Error("panic recovered",
					"request_id", reqID,
					"error", rvr,
					"stack", string(debug.Stack()),
				)

				if err := utils.ErrorJSON(w, http.StatusInternalServerError, utils.CodeInternalError, "internal server error"); err != nil {
					mw.Logger.Error("failed to write panic error response", "request_id", reqID, "error", err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
