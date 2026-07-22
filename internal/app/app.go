package app

import (
	"context"
	"fmt"
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
	cfg, err := getServerConfig()
	if err != nil {
		return fmt.Errorf("get server config: %w", err)
	}

	storage := localCache.New()
	domain := server.New(storage)
	h := handler.New(domain)
	router := router.New(h)

	slog.Info(fmt.Sprintf("SERVER START %s", cfg.Address))
	return http.ListenAndServe(cfg.Address, router)
}

func RunAgent(ctx context.Context) error {
	cfg, err := getAgentConfig()
	if err != nil {
		return fmt.Errorf("get agent config: %w", err)
	}

	storage := localCache.New()
	ag := runtime.New(&runtime.Config{PollInterval: time.Duration(cfg.PollInterval) * time.Second}, storage)
	cl := httpclient.New(&httpclient.Config{
		Timeout: 2 * time.Second,
		URL:     "http://" + cfg.ServerAddress,
	})

	domain := agent.New(&agent.Config{ReportInterval: time.Duration(cfg.ReportInterval) * time.Second}, ag, cl)

	return domain.Start(ctx)
}
