package controller

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

type IController interface {
	SetMetric(ctx context.Context, metric *models.Metrics) error
}
