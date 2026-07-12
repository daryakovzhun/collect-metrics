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

	r.Route("/update", func(r chi.Router) {
		r.Post("/{metric_type}/{metric_name}/{metric_value}", h.SetMetric)
		r.Post("/", h.UpdateMetricFromBody)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{metric_type}/{metric_name}", h.GetMetric)
		r.Post("/", h.GetMetricFromBody)
	})

	r.Get("/", h.GetAllMetrics)

	return r
}
