// internal/middleware/recoverer.go
package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				// http.ErrAbortHandler is a special sentinel panic in net/http;
				// it should not be recovered or logged as an error.
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				reqID := GetRequestID(r.Context())
				log.Printf("[PANIC] request_id=%s err=%v\n%s", reqID, rvr, debug.Stack())

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
