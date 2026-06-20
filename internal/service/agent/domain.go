package agent

import (
	"context"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
	"time"
)

type Config struct {
	ReportInterval time.Duration
}

type domain struct {
	cfg *Config
	repo
}

func New(cfg *Config, repo repository.IRepository) controller.IAgentController {
	return &domain{
		cfg:  cfg,
		repo: repo,
	}
}

func (d *domain) Start(ctx context.Context) error {
	ticker := time.NewTicker(d.cfg.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done!")
			return nil
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
			return nil
		}
	}
}
