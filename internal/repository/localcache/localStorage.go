package localcache

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
	"sync"
)

type storage struct {
	gauge   map[string]models.Metrics
	counter map[string]models.Metrics
	mu      sync.Mutex
}

func New() repository.IRepository {
	return &storage{
		gauge:   make(map[string]models.Metrics),
		counter: make(map[string]models.Metrics),
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return nil
}

func (s *storage) SetGaugeMetric(ctx context.Context, metric models.Metrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gauge[metric.ID] = metric
	return nil
}

func (s *storage) SetCounterMetric(ctx context.Context, metric models.Metrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.counter[metric.ID]
	if !ok {
		s.counter[metric.ID] = metric
		return nil
	}

	delta := utils.FromPointer(val.Delta) + utils.FromPointer(metric.Delta)
	metric.Delta = utils.ToPointer(delta)

	s.counter[metric.ID] = metric
	return nil
}

func (s *storage) GetGaugeMetrics(ctx context.Context) ([]models.Metrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gauge := make([]models.Metrics, 0, len(s.gauge))
	for _, v := range s.gauge {
		gauge = append(gauge, v)
	}

	return gauge, nil
}

func (s *storage) GetCounterMetrics(ctx context.Context) ([]models.Metrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	counter := make([]models.Metrics, 0, len(s.counter))
	for _, v := range s.counter {
		counter = append(counter, v)
	}

	return counter, nil
}

func (s *storage) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics := make([]models.Metrics, 0, len(s.gauge)+len(s.counter))
	for _, v := range s.gauge {
		metrics = append(metrics, v)
	}

	for _, v := range s.counter {
		metrics = append(metrics, v)
	}

	return metrics, nil
}

func (s *storage) GetMetricByID(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
