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
				a.collectGaugeMetrics(memStats)
				a.collectCounterMetrics()
			}
		}
	}()
}

func (a *rtAgent) collectGaugeMetrics(memStats runtime.MemStats) {
	a.SetGaugeMetric(toGaugeMetric(Alloc, float64(memStats.Alloc)))
	a.SetGaugeMetric(toGaugeMetric(BuckHashSys, float64(memStats.BuckHashSys)))
	a.SetGaugeMetric(toGaugeMetric(Frees, float64(memStats.Frees)))
	a.SetGaugeMetric(toGaugeMetric(GCCPUFraction, memStats.GCCPUFraction))
	a.SetGaugeMetric(toGaugeMetric(GCSys, float64(memStats.GCSys)))
	a.SetGaugeMetric(toGaugeMetric(HeapAlloc, float64(memStats.HeapAlloc)))
	a.SetGaugeMetric(toGaugeMetric(HeapIdle, float64(memStats.HeapIdle)))
	a.SetGaugeMetric(toGaugeMetric(HeapInuse, float64(memStats.HeapInuse)))
	a.SetGaugeMetric(toGaugeMetric(HeapObjects, float64(memStats.HeapObjects)))
	a.SetGaugeMetric(toGaugeMetric(HeapReleased, float64(memStats.HeapReleased)))
	a.SetGaugeMetric(toGaugeMetric(HeapSys, float64(memStats.HeapSys)))
	a.SetGaugeMetric(toGaugeMetric(LastGC, float64(memStats.LastGC)))
	a.SetGaugeMetric(toGaugeMetric(Lookups, float64(memStats.Lookups)))
	a.SetGaugeMetric(toGaugeMetric(MCacheInuse, float64(memStats.MCacheInuse)))
	a.SetGaugeMetric(toGaugeMetric(MCacheSys, float64(memStats.MCacheSys)))
	a.SetGaugeMetric(toGaugeMetric(MSpanInuse, float64(memStats.MSpanInuse)))
	a.SetGaugeMetric(toGaugeMetric(MSpanSys, float64(memStats.MSpanSys)))
	a.SetGaugeMetric(toGaugeMetric(Mallocs, float64(memStats.Mallocs)))
	a.SetGaugeMetric(toGaugeMetric(NextGC, float64(memStats.NextGC)))
	a.SetGaugeMetric(toGaugeMetric(NumForcedGC, float64(memStats.NumForcedGC)))
	a.SetGaugeMetric(toGaugeMetric(NumGC, float64(memStats.NumGC)))
	a.SetGaugeMetric(toGaugeMetric(OtherSys, float64(memStats.OtherSys)))
	a.SetGaugeMetric(toGaugeMetric(PauseTotalNs, float64(memStats.PauseTotalNs)))
	a.SetGaugeMetric(toGaugeMetric(StackInuse, float64(memStats.StackInuse)))
	a.SetGaugeMetric(toGaugeMetric(StackSys, float64(memStats.StackSys)))
	a.SetGaugeMetric(toGaugeMetric(Sys, float64(memStats.Sys)))
	a.SetGaugeMetric(toGaugeMetric(TotalAlloc, float64(memStats.TotalAlloc)))
}

func (a *rtAgent) collectCounterMetrics() {
	a.SetCounterMetric(toCounterMetric(PollCount, 1))
	a.SetCounterMetric(toCounterMetric(RandomValue, rand.Int64()))
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
