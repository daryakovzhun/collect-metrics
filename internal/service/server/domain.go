package server

import (
	"context"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	StoreInterval time.Duration
	Restore       bool
}

type domain struct {
	cfg         *Config
	repo        repository.IRepository
	database    repository.IRepository
	fileStorage repository.IFile
}

func New(ctx context.Context, cfg *Config, repo, database repository.IRepository,
	fileStorage repository.IFile) controller.IServerController {
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		database:    database,
		fileStorage: fileStorage,
	}

	if cfg.StoreInterval > 0 {
		go d.asyncSaveMetrics(ctx)
	}

	if cfg.Restore {
		d.restoreMetrics(ctx)
	}

	return d
}

func (d *domain) restoreMetrics(ctx context.Context) {
	metrics, err := d.fileStorage.Read()
	if err != nil {
		logger.Log.Error("failed to read metrics", zap.Error(err))
		return
	}

	for _, m := range metrics {
		if err = d.SetMetric(ctx, &m); err != nil {
			logger.Log.Error("failed to set metric", zap.Error(err))
		}
	}
}

func (d *domain) asyncSaveMetrics(ctx context.Context) {
	ticker := time.NewTicker(d.cfg.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := d.uploadMetrics(); err != nil {
				logger.Log.Error("failed to save metrics", zap.Error(err))
				continue
			}

			logger.Log.Info("upload metrics to file storage")
		}
	}
}

func (d *domain) uploadMetrics() error {
	metrics, err := d.repo.GetAllMetrics()
	if err != nil {
		return fmt.Errorf("failed to get all metrics from localcahe: %w", err)
	}

	err = d.fileStorage.Write(metrics)
	if err != nil {
		return fmt.Errorf("failed to write metrics to file storage: %w", err)
	}

	return nil
}

func (d *domain) SetMetric(ctx context.Context, metric *models.Metrics) error {
	switch metric.MType {
	case models.Gauge:
		d.repo.SetGaugeMetric(*metric)
	case models.Counter:
		d.repo.SetCounterMetric(*metric)
	default:
		return models.ErrUnknownMetricType
	}

	if d.cfg.StoreInterval == 0 {
		err := d.fileStorage.Write([]models.Metrics{*metric})
		if err != nil {
			return fmt.Errorf("failed to write metrics to file storage: %w", err)
		}
	}

	return nil
}

func (d *domain) GetMetric(ctx context.Context, metric *models.Metrics) (models.Metrics, error) {
	return d.repo.GetMetricByID(*metric)
}

func (d *domain) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	return d.repo.GetAllMetrics()
}

func (d *domain) Ping(ctx context.Context) error {
	return d.database.Ping(ctx)
}
