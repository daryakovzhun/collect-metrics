package runtime

import (
	"context"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/agent"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
	"math/rand/v2"
	"runtime"
	"time"
)

type Config struct {
	PollInterval time.Duration
}

type rtAgent struct {
	cfg *Config
	repository.IRepository
}

func New(cfg *Config, repo repository.IRepository) agent.IAgent {
	return &rtAgent{
		cfg:         cfg,
		IRepository: repo,
	}
}

func (a *rtAgent) Collect(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(a.cfg.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				a.collectGaugeMetrics(ctx, memStats)
				a.collectCounterMetrics(ctx)
				err := a.collectSystemMetrics(ctx)
				if err != nil {
					logger.Log.Error("collect system metrics failed", zap.Error(err))
				}
			}
		}
	}()
}

func (a *rtAgent) collectGaugeMetrics(ctx context.Context, memStats runtime.MemStats) {
	a.SetGaugeMetric(ctx, toGaugeMetric(Alloc, float64(memStats.Alloc)))
	a.SetGaugeMetric(ctx, toGaugeMetric(BuckHashSys, float64(memStats.BuckHashSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(Frees, float64(memStats.Frees)))
	a.SetGaugeMetric(ctx, toGaugeMetric(GCCPUFraction, memStats.GCCPUFraction))
	a.SetGaugeMetric(ctx, toGaugeMetric(GCSys, float64(memStats.GCSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(HeapAlloc, float64(memStats.HeapAlloc)))
	a.SetGaugeMetric(ctx, toGaugeMetric(HeapIdle, float64(memStats.HeapIdle)))
	a.SetGaugeMetric(ctx, toGaugeMetric(HeapInuse, float64(memStats.HeapInuse)))
	a.SetGaugeMetric(ctx, toGaugeMetric(HeapObjects, float64(memStats.HeapObjects)))
	a.SetGaugeMetric(ctx, toGaugeMetric(HeapReleased, float64(memStats.HeapReleased)))
	a.SetGaugeMetric(ctx, toGaugeMetric(HeapSys, float64(memStats.HeapSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(LastGC, float64(memStats.LastGC)))
	a.SetGaugeMetric(ctx, toGaugeMetric(Lookups, float64(memStats.Lookups)))
	a.SetGaugeMetric(ctx, toGaugeMetric(MCacheInuse, float64(memStats.MCacheInuse)))
	a.SetGaugeMetric(ctx, toGaugeMetric(MCacheSys, float64(memStats.MCacheSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(MSpanInuse, float64(memStats.MSpanInuse)))
	a.SetGaugeMetric(ctx, toGaugeMetric(MSpanSys, float64(memStats.MSpanSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(Mallocs, float64(memStats.Mallocs)))
	a.SetGaugeMetric(ctx, toGaugeMetric(NextGC, float64(memStats.NextGC)))
	a.SetGaugeMetric(ctx, toGaugeMetric(NumForcedGC, float64(memStats.NumForcedGC)))
	a.SetGaugeMetric(ctx, toGaugeMetric(NumGC, float64(memStats.NumGC)))
	a.SetGaugeMetric(ctx, toGaugeMetric(OtherSys, float64(memStats.OtherSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(PauseTotalNs, float64(memStats.PauseTotalNs)))
	a.SetGaugeMetric(ctx, toGaugeMetric(StackInuse, float64(memStats.StackInuse)))
	a.SetGaugeMetric(ctx, toGaugeMetric(StackSys, float64(memStats.StackSys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(Sys, float64(memStats.Sys)))
	a.SetGaugeMetric(ctx, toGaugeMetric(TotalAlloc, float64(memStats.TotalAlloc)))
	a.SetGaugeMetric(ctx, toGaugeMetric(RandomValue, rand.Float64()))
}

func (a *rtAgent) collectSystemMetrics(ctx context.Context) error {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get virtual memory: %w", err)
	}

	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return fmt.Errorf("failed to get cpu data: %w", err)
	}

	var cpuUtil float64
	if len(cpuPercents) > 0 {
		cpuUtil = cpuPercents[0]
	}

	a.SetGaugeMetric(ctx, toGaugeMetric(TotalMemory, float64(vMem.Total)))
	a.SetGaugeMetric(ctx, toGaugeMetric(FreeMemory, float64(vMem.Free)))
	a.SetGaugeMetric(ctx, toGaugeMetric(CPUutilization1, cpuUtil))

	return nil
}

func (a *rtAgent) collectCounterMetrics(ctx context.Context) {
	a.SetCounterMetric(ctx, toCounterMetric(PollCount, 1))
}

func toGaugeMetric(name string, value float64) models.Metrics {
	return models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: utils.ToPointer(value),
	}
}

func toCounterMetric(name string, value int64) models.Metrics {
	return models.Metrics{
		ID:    name,
		MType: models.Counter,
		Delta: utils.ToPointer(value),
	}
}
