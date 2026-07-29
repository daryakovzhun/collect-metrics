package agent

import (
	"context"
	"time"

	"github.com/daryakovzhun/collect-metrics/internal/agent"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/service/controller"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Config содержит настройки контроллера агента.
type Config struct {
	ReportInterval time.Duration // интервал сбора и отправки метрик
	RateLimit      int           // количество параллельных воркеров для отправки
}

// domain реализует интерфейс controller.IAgentController.
type domain struct {
	cfg    *Config
	agent  agent.IAgent
	client client.IClient
}

// New создаёт новый экземпляр контроллера агента.
func New(cfg *Config, agent agent.IAgent, client client.IClient) controller.IAgentController {
	return &domain{
		cfg:    cfg,
		agent:  agent,
		client: client,
	}
}

// Start запускает цикл сбора и отправки метрик.
// Блокируется до завершения контекста или ошибки, затем ждёт остановки всех горутин.
func (d *domain) Start(ctx context.Context) {
	logger.Log.Info("starting agent controller",
		zap.Duration("report_interval", d.cfg.ReportInterval),
		zap.Int("rate_limit", d.cfg.RateLimit),
	)

	jobs := make(chan []models.Metrics, d.cfg.RateLimit)

	g, gCtx := errgroup.WithContext(ctx)

	d.startWorkers(gCtx, g, jobs)
	d.startCollector(gCtx, g, jobs)

	if err := g.Wait(); err != nil {
		logger.Log.Error("agent controller stopped with error", zap.Error(err))
	} else {
		logger.Log.Info("agent controller stopped gracefully")
	}
}

// startWorkers запускает пул воркеров, которые читают метрики из канала и отправляют их.
func (d *domain) startWorkers(ctx context.Context, g *errgroup.Group, jobs <-chan []models.Metrics) {
	for i := 0; i < d.cfg.RateLimit; i++ {
		workerID := i
		g.Go(func() error {
			logger.Log.Debug("worker started", zap.Int("worker_id", workerID))
			defer logger.Log.Debug("worker stopped", zap.Int("worker_id", workerID))

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case metrics, ok := <-jobs:
					if !ok {
						return nil
					}
					if err := d.client.SendMetrics(ctx, metrics); err != nil {
						logger.Log.Error("failed to send metrics",
							zap.Int("worker_id", workerID),
							zap.Error(err),
						)
					}
				}
			}
		})
	}
}

func (d *domain) startCollector(ctx context.Context, g *errgroup.Group, jobs chan<- []models.Metrics) {
	g.Go(func() error {
		logger.Log.Debug("collector started")
		defer logger.Log.Debug("collector stopped")
		defer close(jobs)

		ticker := time.NewTicker(d.cfg.ReportInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				metrics, err := d.agent.GetAllMetrics(ctx)
				if err != nil {
					logger.Log.Error("failed to get metrics", zap.Error(err))
					continue
				}

				if len(metrics) == 0 {
					continue
				}

				select {
				case jobs <- metrics:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	})
}
