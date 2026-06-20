package router

import (
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"net/http"
)

func New(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", h.SetMetric)

	return mux
}
