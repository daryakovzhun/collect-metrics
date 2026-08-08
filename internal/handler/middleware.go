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
		if h.cfg.Key == "" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get(hashHeader) == "" {
			http.Error(w, "header HashSHA256 not set", http.StatusBadRequest)
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

		wrapper := &responseWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		respBody := wrapper.buf.Bytes()
		hashResp, err := utils.ComputeHash(h.cfg.Key, respBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set(hashHeader, hex.EncodeToString(hashResp))

		if !wrapper.headerWritten {
			w.WriteHeader(wrapper.status)
		}
		if len(respBody) > 0 {
			_, _ = w.Write(respBody)
		}
	})
}

type responseWrapper struct {
	http.ResponseWriter
	buf           bytes.Buffer
	status        int
	headerWritten bool
}

func (rw *responseWrapper) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	return rw.buf.Write(b)
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.headerWritten = true
}
