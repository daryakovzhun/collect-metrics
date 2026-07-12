package app

import (
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"github.com/daryakovzhun/collect-metrics/internal/repository/localCache"
	"github.com/daryakovzhun/collect-metrics/internal/server"
	"github.com/daryakovzhun/collect-metrics/internal/service"
	"log/slog"
	"net/http"
)

func Run() error {
	storage := localCache.New()
	domain := service.New(storage)
	h := handler.New(domain)
	router := server.New(h)

	slog.Info("SERVER START :8080")
	return http.ListenAndServe(":8080", router)
}
