package agent

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

//go:generate mockgen -source=./interface.go -destination=./../mocks/agent.go -package=mocks

type IAgent interface {
	Collect(ctx context.Context)
	GetAllMetrics() ([]models.Metrics, error)
}
