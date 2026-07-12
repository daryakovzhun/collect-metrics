package server

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
)

type domain struct {
	repo repository.IRepository
}

func New(repo repository.IRepository) controller.IServerController {
	return &domain{
		repo: repo,
	}
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

	return nil
}

func (d *domain) GetMetric(ctx context.Context, metric *models.Metrics) (models.Metrics, error) {
	return d.repo.GetMetricByID(*metric)
}

func (d *domain) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	return d.repo.GetAllMetrics()
}
