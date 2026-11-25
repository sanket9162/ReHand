// internal/middleware/recoverer.go
package middleware

import (
	"net/http"
	"runtime/debug"
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

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
