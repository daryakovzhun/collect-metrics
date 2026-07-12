package repository

import models "github.com/daryakovzhun/collect-metrics/internal/model"

//go:generate mockgen -source=./interface.go -destination=./../mocks/repository.go -package=mocks

type IRepository interface {
	SetGaugeMetric(metric models.Metrics)
	SetCounterMetric(metric models.Metrics)
	GetGaugeMetrics() ([]models.Metrics, error)
	GetCounterMetrics() ([]models.Metrics, error)
	GetAllMetrics() ([]models.Metrics, error)
	GetMetricByID(metric models.Metrics) (models.Metrics, error)
}
