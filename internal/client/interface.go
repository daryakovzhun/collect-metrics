package client

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../mocks/client.go -package=mocks

type IClient interface {
	SendMetric(ctx context.Context, metric *models.Metrics) error
	SendMetrics(ctx context.Context, metric []models.Metrics) error
}
