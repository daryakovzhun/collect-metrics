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

func (s *storage) SetGaugeMetric(metric models.Metrics) {
	s.gauge[metric.ID] = metric
}

func (s *storage) SetCounterMetric(metric models.Metrics) {
	val, ok := s.counter[metric.ID]
	if !ok {
		s.counter[metric.ID] = metric
	}

	delta := utils.FromPointer(val.Delta) + utils.FromPointer(metric.Delta)
	metric.Delta = utils.ToPointer(delta)
	s.counter[metric.ID] = metric
}

func (s *storage) GetGaugeMetrics() ([]models.Metrics, error) {
	gauge := make([]models.Metrics, 0, len(s.gauge))
	for _, v := range s.gauge {
		gauge = append(gauge, v)
	}

	return gauge, nil
}

func (s *storage) GetCounterMetrics() ([]models.Metrics, error) {
	counter := make([]models.Metrics, 0, len(s.counter))
	for _, v := range s.counter {
		counter = append(counter, v)
	}

	return counter, nil
}

func (s *storage) GetAllMetrics() ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0, len(s.gauge)+len(s.counter))
	for _, v := range s.gauge {
		metrics = append(metrics, v)
	}

	for _, v := range s.counter {
		metrics = append(metrics, v)
	}

	return metrics, nil
}

func (s *storage) GetMetricByID(metric models.Metrics) (models.Metrics, error) {
	var val models.Metrics
	var ok bool

	switch metric.MType {
	case models.Counter:
		val, ok = s.counter[metric.ID]
	case models.Gauge:
		val, ok = s.gauge[metric.ID]
	default:
		return models.Metrics{}, models.ErrUnknownMetricType
	}

	if !ok {
		return models.Metrics{}, models.ErrNotFound
	}

	return val, nil
}
