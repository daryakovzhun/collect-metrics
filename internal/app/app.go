package app

import (
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"github.com/daryakovzhun/collect-metrics/internal/repository/localCache"
	"github.com/daryakovzhun/collect-metrics/internal/router"
	"github.com/daryakovzhun/collect-metrics/internal/service/server"
	"log/slog"
	"net/http"
)

func Run() error {
	storage := localCache.New()
	domain := server.New(storage)
	h := handler.New(domain)
	router := router.New(h)

	slog.Info("SERVER START :8080")
	return http.ListenAndServe(":8080", router)
}

func RunAgent() error {
	storage := localCache.New()

}
