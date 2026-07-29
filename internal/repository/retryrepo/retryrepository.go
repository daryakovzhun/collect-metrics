package retryrepository

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
)

type retryableRepository struct {
	repo repository.IRepository
}

func New(repo repository.IRepository) repository.IRepository {
	return &retryableRepository{repo: repo}
}

func (r *retryableRepository) SetGaugeMetric(ctx context.Context, metric models.Metrics) error {
	return retry(ctx, func() error {
		return r.repo.SetGaugeMetric(ctx, metric)
	})
}

func (r *retryableRepository) SetCounterMetric(ctx context.Context, metric models.Metrics) error {
	return retry(ctx, func() error {
		return r.repo.SetCounterMetric(ctx, metric)
	})
}

func (r *retryableRepository) GetMetricByID(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	var result models.Metrics
	err := retry(ctx, func() error {
		var err error
		result, err = r.repo.GetMetricByID(ctx, metric)
		return err
	})
	return result, err
}

func (r *retryableRepository) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	var result []models.Metrics
	err := retry(ctx, func() error {
		var err error
		result, err = r.repo.GetAllMetrics(ctx)
		return err
	})
	return result, err
}

func (r *retryableRepository) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	return retry(ctx, func() error {
		return r.repo.UpdateMetrics(ctx, metrics)
	})
}

func (r *retryableRepository) Ping(ctx context.Context) error {
	return retry(ctx, func() error {
		return r.repo.Ping(ctx)
	})
}
