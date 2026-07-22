package handler

import (
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func WithLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			next.ServeHTTP(ww, r)

			logger.Log.Info(
				"got incoming HTTP request",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("duration", time.Since(t1).String()),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
			)
		}

		return http.HandlerFunc(fn)
	}
}

func WithGzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
