package controller

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

type IServerController interface {
	SetMetric(ctx context.Context, metric *models.Metrics) error
}

type IAgentController interface {
	Start(ctx context.Context) error
}
