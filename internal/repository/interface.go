package repository

import models "github.com/daryakovzhun/collect-metrics/internal/model"

type IRepository interface {
	SetGaugeMetric(metric *models.Metrics)
	SetCounterMetric(metric *models.Metrics)
	GetGaugeMetrics() []models.Metrics
	GetCounterMetrics() []models.Metrics
}
