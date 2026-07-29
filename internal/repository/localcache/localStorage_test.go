package localcache

import (
	"context"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"testing"

	"github.com/daryakovzhun/collect-metrics/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	repo := New()
	assert.NotNil(t, repo)
}

func TestStorage_SetGaugeMetric(t *testing.T) {
	repo := New().(*storage)

	metric1 := models.Metrics{
		ID:    "test_gauge_1",
		MType: models.Gauge,
		Value: utils.ToPointer(123.45),
	}
	metric2 := models.Metrics{
		ID:    "test_gauge_2",
		MType: models.Gauge,
		Value: utils.ToPointer(67.89),
	}

	repo.SetGaugeMetric(context.Background(), metric1)
	repo.SetGaugeMetric(context.Background(), metric2)

	// Проверяем, что обе метрики сохранены
	gauges, err := repo.GetGaugeMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, gauges, 2)

	// Проверяем значения по ID (можем получить через GetMetricByID)
	for _, m := range gauges {
		if m.ID == metric1.ID {
			assert.Equal(t, *metric1.Value, *m.Value)
		} else if m.ID == metric2.ID {
			assert.Equal(t, *metric2.Value, *m.Value)
		} else {
			t.Errorf("unexpected metric ID: %s", m.ID)
		}
	}

	// Перезаписываем существующую метрику
	newValue := 999.99
	metric1.Value = utils.ToPointer(newValue)
	repo.SetGaugeMetric(context.Background(), metric1)

	gauges, err = repo.GetGaugeMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, gauges, 2)

	found := false
	for _, m := range gauges {
		if m.ID == metric1.ID {
			assert.Equal(t, newValue, *m.Value)
			found = true
			break
		}
	}
	assert.True(t, found, "metric1 not found after overwrite")
}

func TestStorage_SetCounterMetric(t *testing.T) {
	repo := New().(*storage)

	metric := models.Metrics{
		ID:    "test_counter",
		MType: models.Counter,
		Delta: utils.ToPointer(int64(10)),
	}

	repo.SetCounterMetric(context.Background(), metric)

	// Проверяем через GetMetricByID
	retrieved, err := repo.GetMetricByID(context.Background(), metric)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), *retrieved.Delta)

	// Добавляем ещё 5, должно стать 15
	metric.Delta = utils.ToPointer(int64(5))
	repo.SetCounterMetric(context.Background(), metric)

	retrieved, err = repo.GetMetricByID(context.Background(), metric)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), *retrieved.Delta)

	// Добавляем другую counter метрику
	metric2 := models.Metrics{
		ID:    "test_counter_2",
		MType: models.Counter,
		Delta: utils.ToPointer(int64(7)),
	}
	repo.SetCounterMetric(context.Background(), metric2)

	counters, err := repo.GetCounterMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, counters, 2)

	// Проверяем сумму для первой метрики
	var found bool
	for _, m := range counters {
		if m.ID == "test_counter" {
			assert.Equal(t, int64(15), *m.Delta)
			found = true
		}
	}
	assert.True(t, found, "test_counter not found in counters")
}

func TestStorage_GetGaugeMetrics(t *testing.T) {
	repo := New().(*storage)

	// Пустой список
	gauges, err := repo.GetGaugeMetrics(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, gauges)

	// Добавляем метрики
	metric1 := models.Metrics{ID: "g1", MType: models.Gauge, Value: utils.ToPointer(1.1)}
	metric2 := models.Metrics{ID: "g2", MType: models.Gauge, Value: utils.ToPointer(2.2)}
	repo.SetGaugeMetric(context.Background(), metric1)
	repo.SetGaugeMetric(context.Background(), metric2)

	gauges, err = repo.GetGaugeMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, gauges, 2)

	// Проверяем наличие, порядок не важен
	expectedValues := map[string]float64{"g1": 1.1, "g2": 2.2}
	for _, m := range gauges {
		val, ok := expectedValues[m.ID]
		assert.True(t, ok, "unexpected metric ID: %s", m.ID)
		assert.Equal(t, val, *m.Value)
	}
}

func TestStorage_GetCounterMetrics(t *testing.T) {
	repo := New().(*storage)

	// Пустой список
	counters, err := repo.GetCounterMetrics(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, counters)

	// Добавляем метрики
	metric1 := models.Metrics{ID: "c1", MType: models.Counter, Delta: utils.ToPointer(int64(10))}
	metric2 := models.Metrics{ID: "c2", MType: models.Counter, Delta: utils.ToPointer(int64(20))}
	repo.SetCounterMetric(context.Background(), metric1)
	repo.SetCounterMetric(context.Background(), metric2)

	counters, err = repo.GetCounterMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, counters, 2)

	expectedValues := map[string]int64{"c1": 10, "c2": 20}
	for _, m := range counters {
		val, ok := expectedValues[m.ID]
		assert.True(t, ok, "unexpected metric ID: %s", m.ID)
		assert.Equal(t, val, *m.Delta)
	}
}

func TestStorage_GetAllMetrics(t *testing.T) {
	repo := New().(*storage)

	// Пустой список
	all, err := repo.GetAllMetrics(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, all)

	// Добавляем и gauge и counter
	g1 := models.Metrics{ID: "g1", MType: models.Gauge, Value: utils.ToPointer(1.1)}
	g2 := models.Metrics{ID: "g2", MType: models.Gauge, Value: utils.ToPointer(2.2)}
	c1 := models.Metrics{ID: "c1", MType: models.Counter, Delta: utils.ToPointer(int64(10))}
	c2 := models.Metrics{ID: "c2", MType: models.Counter, Delta: utils.ToPointer(int64(20))}

	repo.SetGaugeMetric(context.Background(), g1)
	repo.SetGaugeMetric(context.Background(), g2)
	repo.SetCounterMetric(context.Background(), c1)
	repo.SetCounterMetric(context.Background(), c2)

	all, err = repo.GetAllMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, all, 4)

	// Проверяем, что все присутствуют
	ids := make(map[string]bool)
	for _, m := range all {
		ids[m.ID] = true
	}
	assert.True(t, ids["g1"])
	assert.True(t, ids["g2"])
	assert.True(t, ids["c1"])
	assert.True(t, ids["c2"])

	// Проверяем типы и значения
	for _, m := range all {
		switch m.ID {
		case "g1":
			assert.Equal(t, models.Gauge, m.MType)
			assert.Equal(t, 1.1, *m.Value)
		case "g2":
			assert.Equal(t, models.Gauge, m.MType)
			assert.Equal(t, 2.2, *m.Value)
		case "c1":
			assert.Equal(t, models.Counter, m.MType)
			assert.Equal(t, int64(10), *m.Delta)
		case "c2":
			assert.Equal(t, models.Counter, m.MType)
			assert.Equal(t, int64(20), *m.Delta)
		}
	}
}

func TestStorage_GetMetricByID(t *testing.T) {
	repo := New().(*storage)

	// Добавляем метрики
	gauge := models.Metrics{ID: "gauge1", MType: models.Gauge, Value: utils.ToPointer(3.14)}
	counter := models.Metrics{ID: "counter1", MType: models.Counter, Delta: utils.ToPointer(int64(42))}
	repo.SetGaugeMetric(context.Background(), gauge)
	repo.SetCounterMetric(context.Background(), counter)

	// Успешный поиск gauge
	retrieved, err := repo.GetMetricByID(context.Background(), models.Metrics{ID: "gauge1", MType: models.Gauge})
	assert.NoError(t, err)
	assert.Equal(t, gauge.ID, retrieved.ID)
	assert.Equal(t, gauge.MType, retrieved.MType)
	assert.Equal(t, *gauge.Value, *retrieved.Value)

	// Успешный поиск counter
	retrieved, err = repo.GetMetricByID(context.Background(), models.Metrics{ID: "counter1", MType: models.Counter})
	assert.NoError(t, err)
	assert.Equal(t, counter.ID, retrieved.ID)
	assert.Equal(t, counter.MType, retrieved.MType)
	assert.Equal(t, *counter.Delta, *retrieved.Delta)

	// Несуществующий ID
	_, err = repo.GetMetricByID(context.Background(), models.Metrics{ID: "unknown", MType: models.Gauge})
	assert.ErrorIs(t, err, models.ErrNotFound)

	_, err = repo.GetMetricByID(context.Background(), models.Metrics{ID: "unknown", MType: models.Counter})
	assert.ErrorIs(t, err, models.ErrNotFound)

	// Неизвестный тип
	_, err = repo.GetMetricByID(context.Background(), models.Metrics{ID: "anything", MType: "unknown_type"})
	assert.ErrorIs(t, err, models.ErrUnknownMetricType)

	// Существующий ID, но неправильный тип (например, gauge ищем как counter)
	_, err = repo.GetMetricByID(context.Background(), models.Metrics{ID: "gauge1", MType: models.Counter})
	assert.ErrorIs(t, err, models.ErrNotFound)

	_, err = repo.GetMetricByID(context.Background(), models.Metrics{ID: "counter1", MType: models.Gauge})
	assert.ErrorIs(t, err, models.ErrNotFound)
}
