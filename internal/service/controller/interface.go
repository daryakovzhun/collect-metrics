package controller

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../../mocks/controllers.go -package=mocks

type IServerController interface {
	SetMetric(ctx context.Context, metric *models.Metrics) error
	GetMetric(ctx context.Context, metric *models.Metrics) (models.Metrics, error)
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
}

type IAgentController interface {
	Start(ctx context.Context) error
}
