package repository

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../mocks/repository.go -package=mocks

type IRepository interface {
	SetGaugeMetric(ctx context.Context, metric models.Metrics) error
	SetCounterMetric(ctx context.Context, metric models.Metrics) error
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
	GetMetricByID(ctx context.Context, metric models.Metrics) (models.Metrics, error)
	Ping(ctx context.Context) error
}

type IFile interface {
	Write(metrics []models.Metrics) error
	Read() ([]models.Metrics, error)
}
