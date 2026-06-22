package agent

import (
	"context"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/agent"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
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
	d.agent.Collect(ctx)

	ticker := time.NewTicker(d.cfg.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := d.sendMetrics(); err != nil {
				return err
			}
		}
	}
}

func (d *domain) sendMetrics() error {
	gauge, err := d.agent.GetGaugeMetrics()
	if err != nil {
		return fmt.Errorf("get gauge: %w", err)
	}

	err = d.sendListMetrics(gauge)
	if err != nil {
		return fmt.Errorf("error gauge: %w", err)
	}

	counter, err := d.agent.GetCounterMetrics()
	if err != nil {
		return fmt.Errorf("get counter: %w", err)
	}

	err = d.sendListMetrics(counter)
	if err != nil {
		return fmt.Errorf("error counter: %w", err)
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
