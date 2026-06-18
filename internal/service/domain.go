package service

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
)

type domain struct {
	repo repository.IRepository
}

func New(repo repository.IRepository) controller.IController {
	return &domain{
		repo: repo,
	}
}

func (d *domain) SetMetric(ctx context.Context, metric *models.Metrics) error {
	switch metric.MType {
	case models.Gauge:
		return d.repo.SetGaugeMetric(metric)
	case models.Counter:
		return d.repo.SetCounterMetric(metric)
	}

	return models.ErrUnknownMetricType
}
