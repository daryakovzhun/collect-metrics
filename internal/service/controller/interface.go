package controller

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../../mocks/controllers.go -package=mocks

type IServerController interface {
	SetMetric(ctx context.Context, metric *models.Metrics) error
	UpdateMetrics(ctx context.Context, metrics []models.Metrics) error
	GetMetric(ctx context.Context, metric *models.Metrics) (models.Metrics, error)
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
	Ping(ctx context.Context) error
}

type IAgentController interface {
	Start(ctx context.Context)
}
