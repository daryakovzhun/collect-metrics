package agent

import (
	"context"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/agent"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	ReportInterval time.Duration
}

type domain struct {
	cfg    *Config
	agent  agent.IAgent
	client client.IClient
}

func New(cfg *Config, agent agent.IAgent, client client.IClient) controller.IAgentController {
	return &domain{
		cfg:    cfg,
		agent:  agent,
		client: client,
	}
}

func (d *domain) Start(ctx context.Context) error {
	logger.Log.Info("starting collect metrics")

	d.agent.Collect(ctx)

	ticker := time.NewTicker(d.cfg.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := d.sendMetrics(ctx); err != nil {
				logger.Log.Error("failed to send metrics", zap.Error(err))
				continue
			}
		}
	}
}

func (d *domain) sendMetrics(ctx context.Context) error {
	metrics, err := d.agent.GetAllMetrics(ctx)
	if err != nil {
		return fmt.Errorf("get metrics: %w", err)
	}

	err = d.sendListMetrics(metrics)
	if err != nil {
		return fmt.Errorf("error gauge: %w", err)
	}

	return nil
}

func (d *domain) sendListMetrics(metrics []models.Metrics) error {
	for _, v := range metrics {
		err := d.client.SendMetric(&v)
		if err != nil {
			return fmt.Errorf("failed to send metric: %s: %w", v.ID, err)
		}
	}

	return nil
}
