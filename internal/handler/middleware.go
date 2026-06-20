package handler

import (
	"log/slog"
	"net/http"
	"time"
)

func WithLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		defer func() {
			slog.Info(
				"got incoming HTTP request",
				slog.String("uri", r.RequestURI),
				slog.String("method", r.Method),
				slog.String("duration", time.Since(t1).String()),
			)
		}()

		next.ServeHTTP(w, r)
	}
}
