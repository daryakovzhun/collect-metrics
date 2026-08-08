package router

import (
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(handler.WithLogger)
	r.Use(handler.WithGzipMiddleware)
	r.Use(h.WithHashMiddleware)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricFromBody)
		r.Post("/{metric_type}/{metric_name}/{metric_value}", h.SetMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricFromBody)
		r.Get("/{metric_type}/{metric_name}", h.GetMetric)
	})

	r.Post("/updates/", h.UpdateMetrics)
	r.Get("/ping", h.Ping)
	r.Get("/", h.GetAllMetrics)

	return r
}
