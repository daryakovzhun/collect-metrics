package handler

import (
	"bytes"
	"encoding/hex"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

const hashHeader = "HashSHA256"

func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func WithGzipMiddleware(next http.Handler) http.Handler {
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

		next.ServeHTTP(ow, r)
	})
}

func (h *Handler) WithHashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.cfg.Key == "" || r.Header.Get(hashHeader) == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		reqHash, err := hex.DecodeString(r.Header.Get(hashHeader))
		if err != nil {
			http.Error(w, "invalid hash format", http.StatusBadRequest)
			return
		}

		currentHash, err := utils.ComputeHash(h.cfg.Key, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !utils.CompareHash(currentHash, reqHash) {
			http.Error(w, "hash mismatch", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
