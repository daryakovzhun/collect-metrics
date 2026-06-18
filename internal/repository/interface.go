package repository

import models "github.com/daryakovzhun/collect-metrics/internal/model"

type IRepository interface {
	SetGaugeMetric(metric *models.Metrics) error
	SetCounterMetric(metric *models.Metrics) error
}
