package app

import (
	"context"
	"github.com/daryakovzhun/collect-metrics/internal/agent/runtime"
	httpclient "github.com/daryakovzhun/collect-metrics/internal/client/http"
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"github.com/daryakovzhun/collect-metrics/internal/repository/localCache"
	"github.com/daryakovzhun/collect-metrics/internal/router"
	"github.com/daryakovzhun/collect-metrics/internal/service/agent"
	"github.com/daryakovzhun/collect-metrics/internal/service/server"
	"log/slog"
	"net/http"

	"time"
)

func Run() error {
	storage := localCache.New()
	domain := server.New(storage)
	h := handler.New(domain)
	router := router.New(h)

	slog.Info("SERVER START :8080")
	return http.ListenAndServe(":8080", router)
}

func RunAgent(ctx context.Context) error {
	storage := localCache.New()
	ag := runtime.New(&runtime.Config{PollInterval: 2 * time.Second}, storage)
	cl := httpclient.New(&httpclient.Config{
		Timeout: 2 * time.Second,
		URL:     "http://localhost:8080",
	})

	domain := agent.New(&agent.Config{ReportInterval: 10 * time.Second}, ag, cl)

	return domain.Start(ctx)
}
