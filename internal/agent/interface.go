package agent

import (
	"context"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
)

//go:generate mockgen -source=./interface.go -destination=./../mocks/agent.go -package=mocks

type IAgent interface {
	Collect(ctx context.Context)
	repository.IRepository
}
