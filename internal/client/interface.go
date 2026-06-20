package client

import models "github.com/daryakovzhun/collect-metrics/internal/model"

type IClient interface {
	SendMetric(metric *models.Metrics) error
}
