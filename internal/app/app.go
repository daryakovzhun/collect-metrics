package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/agent/runtime"
	httpclient "github.com/daryakovzhun/collect-metrics/internal/client/http"
	"github.com/daryakovzhun/collect-metrics/internal/handler"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/repository/filestore"
	"github.com/daryakovzhun/collect-metrics/internal/repository/localcache"
	"github.com/daryakovzhun/collect-metrics/internal/repository/pg"
	retryrepository "github.com/daryakovzhun/collect-metrics/internal/repository/retryrepo"
	"github.com/daryakovzhun/collect-metrics/internal/router"
	"github.com/daryakovzhun/collect-metrics/internal/service/agent"
	"github.com/daryakovzhun/collect-metrics/internal/service/server"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"

	"time"
)

func Run(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := getServerConfig()
	if err != nil {
		return fmt.Errorf("get server config: %w", err)
	}

	var storage repository.IRepository
	if len(cfg.DB) > 0 {
		storage, err = pg.New(ctx, &pg.Config{DatabaseDNS: cfg.DB})
		if err != nil {
			return fmt.Errorf("failed to connect database, err: %w", err)
		}
		storage = retryrepository.New(storage)
	} else {
		storage = localcache.New()
	}

	fileStorage := filestore.New(&filestore.Config{Path: cfg.FileStoragePath})

	domain := server.New(egCtx, &server.Config{
		StoreInterval: time.Duration(utils.FromPointer(cfg.StoreInterval)) * time.Second,
		Restore:       utils.FromPointer(cfg.Restore),
	}, storage, fileStorage)
	h := handler.New(domain)
	router := router.New(h)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	eg.Go(func() error {
		logger.Log.Info(fmt.Sprintf("SERVER START %s", cfg.Address))
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(shCtx)
	})

	return eg.Wait()
}

func RunAgent(ctx context.Context) error {
	cfg, err := getAgentConfig()
	if err != nil {
		return fmt.Errorf("get agent config: %w", err)
	}

	storage := localcache.New()
	ag := runtime.New(&runtime.Config{PollInterval: time.Duration(cfg.PollInterval) * time.Second}, storage)
	cl := httpclient.New(&httpclient.Config{
		Timeout: 2 * time.Second,
		URL:     "http://" + cfg.ServerAddress,
	})

	domain := agent.New(&agent.Config{ReportInterval: time.Duration(cfg.ReportInterval) * time.Second}, ag, cl)

	return domain.Start(ctx)
}
