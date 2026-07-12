package router

import (
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(handler.WithLogger())

	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", h.SetMetric)
	r.Get("/value/{metric_type}/{metric_name}", h.GetMetric)
	r.Get("/", h.GetAllMetrics)

	return r
}
