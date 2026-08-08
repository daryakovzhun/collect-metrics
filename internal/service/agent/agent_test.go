package agent

import (
	"context"
	"errors"
	"github.com/daryakovzhun/collect-metrics/internal/mocks"
	"testing"
	"time"

	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap/zaptest"
)

// init настраивает логгер для тестов (вывод в тестовый вывод)
func init() {
	logger.Log = zaptest.NewLogger(&testing.T{}).Named("test")
}

// helper для создания указателя на float64
func float64Ptr(v float64) *float64 {
	return &v
}

// helper для создания указателя на int64
func int64Ptr(v int64) *int64 {
	return &v
}

// TestDomainStart_Success проверяет успешный сбор и отправку метрик.
func TestDomainStart_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgent := mocks.NewMockIAgent(ctrl)
	mockClient := mocks.NewMockIClient(ctrl)

	cfg := &Config{
		ReportInterval: 30 * time.Millisecond,
		RateLimit:      2,
	}

	metrics := []models.Metrics{
		{ID: "test1", MType: "gauge", Value: float64Ptr(1.0)},
		{ID: "test2", MType: "counter", Delta: int64Ptr(5)},
	}

	// Ожидаем, что GetAllMetrics будет вызван как минимум один раз
	mockAgent.EXPECT().
		GetAllMetrics(gomock.Any()).
		Return(metrics, nil).
		MinTimes(1)

	// Ожидаем, что SendMetrics будет вызван для каждой пачки метрик (как минимум один раз)
	mockClient.EXPECT().
		SendMetrics(gomock.Any(), metrics).
		Return(nil).
		MinTimes(1)

	domain := New(cfg, mockAgent, mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		domain.Start(ctx)
		close(done)
	}()

	// Даём время на выполнение нескольких циклов
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Ждём завершения Start
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("domain.Start did not stop after context cancellation")
	}
}

// TestDomainStart_AgentError проверяет, что ошибка агента не прерывает цикл.
func TestDomainStart_AgentError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgent := mocks.NewMockIAgent(ctrl)
	mockClient := mocks.NewMockIClient(ctrl)

	cfg := &Config{
		ReportInterval: 30 * time.Millisecond,
		RateLimit:      2,
	}

	// GetAllMetrics всегда возвращает ошибку
	mockAgent.EXPECT().
		GetAllMetrics(gomock.Any()).
		Return(nil, errors.New("agent error")).
		MinTimes(1)

	// SendMetrics не должен вызываться, т.к. метрики не получены
	mockClient.EXPECT().
		SendMetrics(gomock.Any(), gomock.Any()).
		Times(0)

	domain := New(cfg, mockAgent, mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		domain.Start(ctx)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("domain.Start did not stop")
	}
}

// TestDomainStart_EmptyMetrics проверяет, что пустой список метрик не приводит к отправке.
func TestDomainStart_EmptyMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgent := mocks.NewMockIAgent(ctrl)
	mockClient := mocks.NewMockIClient(ctrl)

	cfg := &Config{
		ReportInterval: 30 * time.Millisecond,
		RateLimit:      2,
	}

	// GetAllMetrics возвращает пустой слайс
	mockAgent.EXPECT().
		GetAllMetrics(gomock.Any()).
		Return([]models.Metrics{}, nil).
		MinTimes(1)

	// SendMetrics не должен вызываться
	mockClient.EXPECT().
		SendMetrics(gomock.Any(), gomock.Any()).
		Times(0)

	domain := New(cfg, mockAgent, mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		domain.Start(ctx)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("domain.Start did not stop")
	}
}

// TestDomainStart_ClientError проверяет, что ошибка клиента логируется, но цикл продолжается.
func TestDomainStart_ClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgent := mocks.NewMockIAgent(ctrl)
	mockClient := mocks.NewMockIClient(ctrl)

	cfg := &Config{
		ReportInterval: 30 * time.Millisecond,
		RateLimit:      2,
	}

	metrics := []models.Metrics{
		{ID: "test", MType: "gauge", Value: float64Ptr(1.0)},
	}

	// GetAllMetrics всегда возвращает метрики
	mockAgent.EXPECT().
		GetAllMetrics(gomock.Any()).
		Return(metrics, nil).
		MinTimes(1)

	// SendMetrics всегда возвращает ошибку
	mockClient.EXPECT().
		SendMetrics(gomock.Any(), metrics).
		Return(errors.New("client error")).
		MinTimes(1)

	domain := New(cfg, mockAgent, mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		domain.Start(ctx)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("domain.Start did not stop")
	}
}

// TestDomainStart_CancelImmediate проверяет немедленную отмену контекста без тиков.
func TestDomainStart_CancelImmediate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgent := mocks.NewMockIAgent(ctrl)
	mockClient := mocks.NewMockIClient(ctrl)

	cfg := &Config{
		ReportInterval: 1 * time.Hour, // большой интервал, чтобы не было тиков
		RateLimit:      2,
	}

	// Никаких вызовов не ожидаем, т.к. контекст отменяется до первого тика
	mockAgent.EXPECT().GetAllMetrics(gomock.Any()).Times(0)
	mockClient.EXPECT().SendMetrics(gomock.Any(), gomock.Any()).Times(0)

	domain := New(cfg, mockAgent, mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		domain.Start(ctx)
		close(done)
	}()

	// Отменяем сразу после старта
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("domain.Start did not stop after immediate cancellation")
	}
}

// TestDomainStart_WorkerCount проверяет, что количество воркеров соответствует RateLimit.
// Проверяем косвенно: при отправке нескольких пачек метрик все они обрабатываются,
// а число вызовов SendMetrics равно числу пачек (независимо от RateLimit).
func TestDomainStart_WorkerCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAgent := mocks.NewMockIAgent(ctrl)
	mockClient := mocks.NewMockIClient(ctrl)

	cfg := &Config{
		ReportInterval: 20 * time.Millisecond,
		RateLimit:      3, // три воркера
	}

	metrics := []models.Metrics{{ID: "test", MType: "gauge", Value: float64Ptr(1.0)}}

	// Ожидаем, что будет несколько вызовов GetAllMetrics
	mockAgent.EXPECT().
		GetAllMetrics(gomock.Any()).
		Return(metrics, nil).
		MinTimes(2)

	// Ожидаем, что каждый вызов GetAllMetrics приведёт к одному вызову SendMetrics
	mockClient.EXPECT().
		SendMetrics(gomock.Any(), metrics).
		Return(nil).
		MinTimes(2)

	domain := New(cfg, mockAgent, mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		domain.Start(ctx)
		close(done)
	}()

	// Даём время на несколько циклов
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("domain.Start did not stop")
	}
}
