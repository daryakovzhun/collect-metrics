package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daryakovzhun/collect-metrics/internal/agent"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	"github.com/daryakovzhun/collect-metrics/internal/mocks"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/golang/mock/gomock"
)

func TestDomain_Start(t *testing.T) {
	type fields struct {
		cfg    *Config
		agent  agent.IAgent
		client client.IClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient)
		wantErr bool
	}{
		{
			name: "successful start and stop by context",
			fields: fields{
				cfg: &Config{ReportInterval: 10 * time.Millisecond},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(20 * time.Millisecond) // даём время на один тик
						cancel()
					}()
					return ctx
				}(),
			},
			setup: func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient) {
				agentMock.EXPECT().Collect(gomock.Any()).Times(1)
				agentMock.EXPECT().GetAllMetrics().Return([]models.Metrics{}, nil).AnyTimes()
				// SendMetric не вызывается, т.к. списки пусты
			},
			wantErr: false,
		},
		{
			name: "error during sendMetrics",
			fields: fields{
				cfg: &Config{ReportInterval: 10 * time.Millisecond},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					// отменим через 100 мс, чтобы успеть получить ошибку
					go func() {
						time.Sleep(100 * time.Millisecond)
						cancel()
					}()
					return ctx
				}(),
			},
			setup: func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient) {
				agentMock.EXPECT().Collect(gomock.Any()).Times(1)
				// Возвращаем метрику, которая вызовет ошибку при отправке
				agentMock.EXPECT().GetAllMetrics().Return([]models.Metrics{
					{ID: "test", Value: toPtrFloat64(1.0)},
				}, errors.New("test")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentMock := mocks.NewMockIAgent(ctrl)
			clientMock := mocks.NewMockIClient(ctrl)
			if tt.setup != nil {
				tt.setup(agentMock, clientMock)
			}
			tt.fields.agent = agentMock
			tt.fields.client = clientMock

			d := &domain{
				cfg:    tt.fields.cfg,
				agent:  tt.fields.agent,
				client: tt.fields.client,
			}

			err := d.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomain_sendMetrics(t *testing.T) {
	type fields struct {
		agent  agent.IAgent
		client client.IClient
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient)
		wantErr bool
	}{
		{
			name: "successful send all metrics",
			setup: func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient) {
				agentMock.EXPECT().GetAllMetrics().Return([]models.Metrics{
					{ID: "g1", Value: toPtrFloat64(1.1)},
					{ID: "g2", Value: toPtrFloat64(2.2)},
					{ID: "c1", Delta: toPtrInt64(10)},
				}, nil).Times(1)

				clientMock.EXPECT().SendMetric(gomock.Any()).Times(3).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error on gauge send",
			setup: func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient) {
				agentMock.EXPECT().GetAllMetrics().Return([]models.Metrics{
					{ID: "g1", Value: toPtrFloat64(1.1)},
				}, nil).Times(1)
				clientMock.EXPECT().SendMetric(gomock.Any()).Return(errors.New("gauge send error")).Times(1)
				// GetCounterMetrics не должен вызываться
			},
			wantErr: true,
		},
		{
			name: "error on counter send",
			setup: func(agentMock *mocks.MockIAgent, clientMock *mocks.MockIClient) {
				agentMock.EXPECT().GetAllMetrics().Return([]models.Metrics{
					{ID: "c1", Delta: toPtrInt64(5)},
				}, nil).Times(1)
				clientMock.EXPECT().SendMetric(gomock.Any()).Return(errors.New("counter send error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			agentMock := mocks.NewMockIAgent(ctrl)
			clientMock := mocks.NewMockIClient(ctrl)
			if tt.setup != nil {
				tt.setup(agentMock, clientMock)
			}
			d := &domain{
				cfg:    &Config{ReportInterval: time.Second}, // не используется в этом тесте
				agent:  agentMock,
				client: clientMock,
			}

			err := d.sendMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("sendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomain_sendListMetrics(t *testing.T) {
	type fields struct {
		client client.IClient
	}
	type args struct {
		metrics []models.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(clientMock *mocks.MockIClient)
		wantErr bool
	}{
		{
			name: "empty list",
			args: args{metrics: []models.Metrics{}},
			setup: func(clientMock *mocks.MockIClient) {
				// никаких вызовов не ожидается
			},
			wantErr: false,
		},
		{
			name: "successful send all",
			args: args{metrics: []models.Metrics{
				{ID: "m1", Value: toPtrFloat64(1.0)},
				{ID: "m2", Delta: toPtrInt64(2)},
			}},
			setup: func(clientMock *mocks.MockIClient) {
				clientMock.EXPECT().SendMetric(gomock.Any()).Return(nil).Times(2)
			},
			wantErr: false,
		},
		{
			name: "error on second metric",
			args: args{metrics: []models.Metrics{
				{ID: "m1", Value: toPtrFloat64(1.0)},
				{ID: "m2", Delta: toPtrInt64(2)},
			}},
			setup: func(clientMock *mocks.MockIClient) {
				gomock.InOrder(
					clientMock.EXPECT().SendMetric(gomock.Any()).Return(nil),
					clientMock.EXPECT().SendMetric(gomock.Any()).Return(errors.New("send failed")),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			clientMock := mocks.NewMockIClient(ctrl)
			if tt.setup != nil {
				tt.setup(clientMock)
			}
			d := &domain{
				client: clientMock,
			}

			err := d.sendListMetrics(tt.args.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendListMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Вспомогательные функции для создания указателей (скопированы из тестов локального кеша)
func toPtrInt64(v int64) *int64 {
	return &v
}

func toPtrFloat64(v float64) *float64 {
	return &v
}
