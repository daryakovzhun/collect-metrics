package localCache

import (
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
)

type storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func New() repository.IRepository {
	return &storage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s *storage) SetGaugeMetric(metric *models.Metrics) error {
	s.gauge[metric.ID] = fromPointer(metric.Value)
	return nil
}

func (s *storage) SetCounterMetric(metric *models.Metrics) error {
	s.counter[metric.ID] += fromPointer(metric.Delta)
	return nil
}

func fromPointer[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
