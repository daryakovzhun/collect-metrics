package localCache

import (
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
)

type storage struct {
	gauge   map[string]models.Metrics
	counter map[string]models.Metrics
}

func New() repository.IRepository {
	return &storage{
		gauge:   make(map[string]models.Metrics),
		counter: make(map[string]models.Metrics),
	}
}

func (s *storage) SetGaugeMetric(metric *models.Metrics) {
	s.gauge[metric.ID] = utils.FromPointer(metric)
}

func (s *storage) SetCounterMetric(metric *models.Metrics) {
	val, ok := s.counter[metric.ID]
	if !ok {
		s.counter[metric.ID] = *metric
	}

	delta := utils.FromPointer(val.Delta) + utils.FromPointer(metric.Delta)
	metric.Delta = utils.ToPointer(delta)
	s.counter[metric.ID] = utils.FromPointer(metric)
}

func (s *storage) GetGaugeMetrics() []models.Metrics {
	gauge := make([]models.Metrics, 0, len(s.gauge))
	for _, v := range s.gauge {
		gauge = append(gauge, v)
	}

	return gauge
}

func (s *storage) GetCounterMetrics() []models.Metrics {
	counter := make([]models.Metrics, 0, len(s.counter))
	for _, v := range s.counter {
		counter = append(counter, v)
	}

	return counter
}
