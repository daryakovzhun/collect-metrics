package runtime

import (
	"context"
	"github.com/daryakovzhun/collect-metrics/internal/agent"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
	"math/rand/v2"
	"runtime"
	"time"
)

type Config struct {
	PollInterval time.Duration
}

type rtAgent struct {
	cfg  *Config
	repo repository.IRepository
}

func New(cfg *Config, repo repository.IRepository) agent.IAgent {
	return &rtAgent{
		cfg:  cfg,
		repo: repo,
	}
}

func (a *rtAgent) Collect(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(a.cfg.PollInterval)
		defer ticker.Stop()

		poolCount := int64(0)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			a.collectGaugeMetrics(memStats)
			a.collectCounterMetrics(poolCount)
			poolCount++
		}
	}()
}

func (a *rtAgent) collectGaugeMetrics(memStats runtime.MemStats) {
	a.repo.SetGaugeMetric(toGaugeMetric(Alloc, float64(memStats.Alloc)))
	a.repo.SetGaugeMetric(toGaugeMetric(BuckHashSys, float64(memStats.BuckHashSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(Frees, float64(memStats.Frees)))
	a.repo.SetGaugeMetric(toGaugeMetric(GCCPUFraction, memStats.GCCPUFraction))
	a.repo.SetGaugeMetric(toGaugeMetric(GCSys, float64(memStats.GCSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(HeapAlloc, float64(memStats.HeapAlloc)))
	a.repo.SetGaugeMetric(toGaugeMetric(HeapIdle, float64(memStats.HeapIdle)))
	a.repo.SetGaugeMetric(toGaugeMetric(HeapInuse, float64(memStats.HeapInuse)))
	a.repo.SetGaugeMetric(toGaugeMetric(HeapObjects, float64(memStats.HeapObjects)))
	a.repo.SetGaugeMetric(toGaugeMetric(HeapReleased, float64(memStats.HeapReleased)))
	a.repo.SetGaugeMetric(toGaugeMetric(HeapSys, float64(memStats.HeapSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(LastGC, float64(memStats.LastGC)))
	a.repo.SetGaugeMetric(toGaugeMetric(Lookups, float64(memStats.Lookups)))
	a.repo.SetGaugeMetric(toGaugeMetric(MCacheInuse, float64(memStats.MCacheInuse)))
	a.repo.SetGaugeMetric(toGaugeMetric(MCacheSys, float64(memStats.MCacheSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(MSpanInuse, float64(memStats.MSpanInuse)))
	a.repo.SetGaugeMetric(toGaugeMetric(MSpanSys, float64(memStats.MSpanSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(Mallocs, float64(memStats.Mallocs)))
	a.repo.SetGaugeMetric(toGaugeMetric(NextGC, float64(memStats.NextGC)))
	a.repo.SetGaugeMetric(toGaugeMetric(NumForcedGC, float64(memStats.NumForcedGC)))
	a.repo.SetGaugeMetric(toGaugeMetric(NumGC, float64(memStats.NumGC)))
	a.repo.SetGaugeMetric(toGaugeMetric(OtherSys, float64(memStats.OtherSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(PauseTotalNs, float64(memStats.PauseTotalNs)))
	a.repo.SetGaugeMetric(toGaugeMetric(StackInuse, float64(memStats.StackInuse)))
	a.repo.SetGaugeMetric(toGaugeMetric(StackSys, float64(memStats.StackSys)))
	a.repo.SetGaugeMetric(toGaugeMetric(Sys, float64(memStats.Sys)))
	a.repo.SetGaugeMetric(toGaugeMetric(TotalAlloc, float64(memStats.TotalAlloc)))
}

func (a *rtAgent) collectCounterMetrics(poolCount int64) {
	a.repo.SetCounterMetric(toCounterMetric(PollCount, poolCount))
	a.repo.SetCounterMetric(toCounterMetric(RandomValue, rand.Int64()))
}

func toGaugeMetric(name string, value float64) *models.Metrics {
	return &models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: utils.ToPointer(value),
	}
}

func toCounterMetric(name string, value int64) *models.Metrics {
	return &models.Metrics{
		ID:    name,
		MType: models.Counter,
		Delta: utils.ToPointer(value),
	}
}
