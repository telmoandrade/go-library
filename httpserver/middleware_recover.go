package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// MiddlewareRecover is a middleware that recovers from panics, logs the panic, and responds with an HTTP status of 500 (Internal Server Error).
func MiddlewareRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				slog.ErrorContext(r.Context(), "Panic recover", slog.Group("error",
					slog.Bool("recover", true),
					slog.String("message", fmt.Sprintf("%v", err)),
					slog.String("stack", string(debug.Stack())),
				))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
