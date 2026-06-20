package agent

import (
	"context"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
)

type IAgent interface {
	Collect(ctx context.Context)
	repository.IRepository
}
