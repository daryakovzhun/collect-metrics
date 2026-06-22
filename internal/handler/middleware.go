package handler

import (
	"log/slog"
	"net/http"
	"time"
)

func WithLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		defer func() {
			slog.Info(
				"got incoming HTTP request",
				slog.String("uri", r.RequestURI),
				slog.String("method", r.Method),
				slog.String("duration", time.Since(t1).String()),
			)
		}()

		h.ServeHTTP(w, r)
	})
}
