package localCache

import (
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
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

func (s *storage) SetGaugeMetric(metric *models.Metrics) {
	s.gauge[metric.ID] = utils.FromPointer(metric.Value)
}

func (s *storage) SetCounterMetric(metric *models.Metrics) {
	s.counter[metric.ID] += utils.FromPointer(metric.Delta)
}
