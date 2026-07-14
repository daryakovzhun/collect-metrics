package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daryakovzhun/collect-metrics/internal/mocks"
	"github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNew_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	metrics := []models.Metrics{
		{ID: "metric1", MType: models.Gauge, Value: func() *float64 { v := 123.45; return &v }()},
		{ID: "metric2", MType: models.Counter, Delta: func() *int64 { v := int64(10); return &v }()},
	}
	file.EXPECT().Read().Return(metrics, nil)

	repo.EXPECT().SetGaugeMetric(gomock.Any()).Times(1)
	repo.EXPECT().SetCounterMetric(gomock.Any()).Times(1)

	cfg := &Config{
		StoreInterval: 5, // или >0, если нужно проверить асинхронное сохранение
		Restore:       true,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	d.restoreMetrics(ctx)
}

func TestNew_Restore_ReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	file.EXPECT().Read().Return(nil, errors.New("read error"))
	// Никаких вызовов SetMetric не ожидается
	repo.EXPECT().SetGaugeMetric(gomock.Any()).Times(0)
	repo.EXPECT().SetCounterMetric(gomock.Any()).Times(0)

	cfg := &Config{
		StoreInterval: 0,
		Restore:       true,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}
	d.restoreMetrics(ctx)
}

func TestNew_NoRestore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	file.EXPECT().Read().Times(1)

	cfg := &Config{
		StoreInterval: 0,
		Restore:       false,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}
	d.restoreMetrics(ctx)
}

func TestSetMetric_Gauge_StoreIntervalZero(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 0}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	metric := &models.Metrics{
		ID:    "test_gauge",
		MType: models.Gauge,
		Value: func() *float64 { v := 99.9; return &v }(),
	}

	repo.EXPECT().SetGaugeMetric(*metric).Times(1)
	file.EXPECT().Write([]models.Metrics{*metric}).Return(nil).Times(1)

	err := d.SetMetric(ctx, metric)
	assert.NoError(t, err)
}

func TestSetMetric_Counter_StoreIntervalZero(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 0}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	metric := &models.Metrics{
		ID:    "test_counter",
		MType: models.Counter,
		Delta: func() *int64 { v := int64(5); return &v }(),
	}

	repo.EXPECT().SetCounterMetric(*metric).Times(1)
	file.EXPECT().Write([]models.Metrics{*metric}).Return(nil).Times(1)

	err := d.SetMetric(ctx, metric)
	assert.NoError(t, err)
}

func TestSetMetric_UnknownType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 0}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	metric := &models.Metrics{
		ID:    "unknown",
		MType: "unknown",
	}

	repo.EXPECT().SetGaugeMetric(gomock.Any()).Times(0)
	repo.EXPECT().SetCounterMetric(gomock.Any()).Times(0)
	file.EXPECT().Write(gomock.Any()).Times(0)

	err := d.SetMetric(ctx, metric)
	assert.ErrorIs(t, err, models.ErrUnknownMetricType)
}

func TestSetMetric_StoreIntervalPositive_NoSyncWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 10 * time.Second}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	metric := &models.Metrics{
		ID:    "test_gauge",
		MType: models.Gauge,
		Value: func() *float64 { v := 1.23; return &v }(),
	}

	repo.EXPECT().SetGaugeMetric(*metric).Times(1)
	file.EXPECT().Write(gomock.Any()).Times(0)

	err := d.SetMetric(ctx, metric)
	assert.NoError(t, err)
}

//func TestGetMetric(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	repo := mocks.NewMockIRepository(ctrl)
//	file := mocks.NewMockIFile(ctrl)
//
//	cfg := &Config{StoreInterval: 0}
//	ctx := context.Background()
//	d := &domain{
//		cfg:         cfg,
//		repo:        repo,
//		fileStorage: file,
//	}
//
//	input := &models.Metrics{ID: "metric1", MType: models.Gauge}
//	expected := models.Metrics{ID: "metric1", MType: models.Gauge, Value: func() *float64 { v := 42.0; return &v }()}
//
//	repo.EXPECT().GetMetricByID(*input).Return(expected, nil)
//
//	result, err := d.GetMetric(ctx, input)
//	assert.NoError(t, err)
//	assert.Equal(t, expected, result)
//}

func TestGetMetric_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 0}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	input := &models.Metrics{ID: "unknown", MType: models.Gauge}
	repo.EXPECT().GetMetricByID(*input).Return(models.Metrics{}, models.ErrNotFound)

	_, err := d.GetMetric(ctx, input)
	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestGetAllMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 0}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	expected := []models.Metrics{
		{ID: "m1", MType: models.Gauge, Value: func() *float64 { v := 1.1; return &v }()},
		{ID: "m2", MType: models.Counter, Delta: func() *int64 { v := int64(2); return &v }()},
	}
	repo.EXPECT().GetAllMetrics().Return(expected, nil)

	result, err := d.GetAllMetrics(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetAllMetrics_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIRepository(ctrl)
	file := mocks.NewMockIFile(ctrl)

	cfg := &Config{StoreInterval: 0}
	ctx := context.Background()
	d := &domain{
		cfg:         cfg,
		repo:        repo,
		fileStorage: file,
	}

	repo.EXPECT().GetAllMetrics().Return(nil, errors.New("db error"))

	_, err := d.GetAllMetrics(ctx)
	assert.Error(t, err)
}
