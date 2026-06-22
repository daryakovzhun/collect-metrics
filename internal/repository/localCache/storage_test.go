package localCache

import (
	"testing"

	models "github.com/daryakovzhun/collect-metrics/internal/model"
)

// Вспомогательные функции для создания указателей
func toPtrInt64(v int64) *int64 {
	return &v
}

func toPtrFloat64(v float64) *float64 {
	return &v
}

// Тест для New()
func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Error("New() returned nil")
	}
	store, ok := s.(*storage)
	if !ok {
		t.Fatal("New() returned not *storage")
	}
	if store.gauge == nil {
		t.Error("gauge map is nil")
	}
	if store.counter == nil {
		t.Error("counter map is nil")
	}
}

// Тест SetGaugeMetric – добавление и проверка копирования
func TestStorage_SetGaugeMetric(t *testing.T) {
	s := New().(*storage)
	metric := &models.Metrics{
		ID:    "test_gauge",
		Value: toPtrFloat64(123.45),
	}
	s.SetGaugeMetric(metric)

	if len(s.gauge) != 1 {
		t.Errorf("expected 1 gauge metric, got %d", len(s.gauge))
	}
	val, ok := s.gauge["test_gauge"]
	if !ok {
		t.Error("metric not found in gauge map")
	}
	if val.Value == nil {
		t.Error("Value is nil")
	} else if *val.Value != 123.45 {
		t.Errorf("expected value 123.45, got %v", *val.Value)
	}

	// Проверка, что сохранённое значение не зависит от внешнего изменения
	metric.Value = toPtrFloat64(999.99)
	if *s.gauge["test_gauge"].Value == 999.99 {
		t.Error("saved metric was modified by external change")
	}
}

// Тест SetGaugeMetric – перезапись существующей метрики
func TestStorage_SetGaugeMetricOverwrite(t *testing.T) {
	s := New().(*storage)
	s.SetGaugeMetric(&models.Metrics{ID: "g", Value: toPtrFloat64(1.0)})
	s.SetGaugeMetric(&models.Metrics{ID: "g", Value: toPtrFloat64(2.0)})
	if *s.gauge["g"].Value != 2.0 {
		t.Errorf("expected 2.0, got %v", *s.gauge["g"].Value)
	}
}

// Тест SetCounterMetric – добавление и суммирование
func TestStorage_SetCounterMetric(t *testing.T) {
	s := New().(*storage)

	// Первое добавление
	metric1 := &models.Metrics{
		ID:    "test_counter",
		Delta: toPtrInt64(10),
	}
	s.SetCounterMetric(metric1)

	if len(s.counter) != 1 {
		t.Errorf("expected 1 counter metric, got %d", len(s.counter))
	}
	val, ok := s.counter["test_counter"]
	if !ok {
		t.Error("metric not found in counter map")
	}
	if val.Delta == nil {
		t.Error("Delta is nil")
	} else if *val.Delta != 10 {
		t.Errorf("expected Delta 10, got %d", *val.Delta)
	}

	// Второе добавление – должно суммироваться
	metric2 := &models.Metrics{
		ID:    "test_counter",
		Delta: toPtrInt64(5),
	}
	s.SetCounterMetric(metric2)

	val2, ok := s.counter["test_counter"]
	if !ok {
		t.Error("metric not found after second add")
	}
	if *val2.Delta != 15 {
		t.Errorf("expected Delta 15 after sum, got %d", *val2.Delta)
	}

	// Проверка, что переданный metric был изменён (на сумму)
	if *metric2.Delta != 15 {
		t.Errorf("metric2.Delta was modified to %d, expected 15", *metric2.Delta)
	}
}

// Тест SetCounterMetric – обработка nil Delta
func TestStorage_SetCounterMetricWithNilDelta(t *testing.T) {
	s := New().(*storage)
	metric := &models.Metrics{ID: "nil_delta", Delta: nil}
	s.SetCounterMetric(metric)

	val, ok := s.counter["nil_delta"]
	if !ok {
		t.Error("metric not added")
	}
	if val.Delta == nil {
		t.Error("Delta is nil, expected 0")
	} else if *val.Delta != 0 {
		t.Errorf("expected Delta 0, got %d", *val.Delta)
	}
	// Проверка, что переданный метрик обновился
	if metric.Delta == nil {
		t.Error("metric.Delta should be updated to pointer to 0")
	} else if *metric.Delta != 0 {
		t.Errorf("metric.Delta expected 0, got %d", *metric.Delta)
	}
}

// Тест GetGaugeMetrics – возврат всех gauge метрик
func TestStorage_GetGaugeMetrics(t *testing.T) {
	s := New().(*storage)
	s.SetGaugeMetric(&models.Metrics{ID: "g1", Value: toPtrFloat64(1.1)})
	s.SetGaugeMetric(&models.Metrics{ID: "g2", Value: toPtrFloat64(2.2)})

	metrics, _ := s.GetGaugeMetrics()
	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}
	found := make(map[string]bool)
	for _, m := range metrics {
		found[m.ID] = true
	}
	if !found["g1"] || !found["g2"] {
		t.Errorf("missing expected IDs: got %v", found)
	}
}

// Тест GetCounterMetrics – возврат всех counter метрик
func TestStorage_GetCounterMetrics(t *testing.T) {
	s := New().(*storage)
	s.SetCounterMetric(&models.Metrics{ID: "c1", Delta: toPtrInt64(10)})
	s.SetCounterMetric(&models.Metrics{ID: "c2", Delta: toPtrInt64(20)})

	metrics, _ := s.GetCounterMetrics()
	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}
	found := make(map[string]bool)
	for _, m := range metrics {
		found[m.ID] = true
	}
	if !found["c1"] || !found["c2"] {
		t.Errorf("missing expected IDs: got %v", found)
	}
}
