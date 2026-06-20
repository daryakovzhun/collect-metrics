package client

import models "github.com/daryakovzhun/collect-metrics/internal/model"

//go:generate mockgen -source=./interface.go -destination=./../mocks/client.go -package=mocks

type IClient interface {
	SendMetric(metric *models.Metrics) error
}
