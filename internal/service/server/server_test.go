package server

import (
	"github.com/daryakovzhun/collect-metrics/internal/mocks"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockIRepository(ctrl)

	s := New(store)

	assert.NotEmpty(t, s)
}

func TestService_UpdateCounterMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		metric  *models.Metrics
		prepare func(store *mocks.MockIRepository, metric *models.Metrics)
		error   error
	}{
		{
			name: "updating counter metric",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIRepository, metric *models.Metrics) {
				store.EXPECT().SetCounterMetric(metric)
			},
		},
		{
			name: "updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIRepository, metric *models.Metrics) {
				store.EXPECT().SetGaugeMetric(metric)
			},
		}, {
			name:   "updating unknown metric",
			metric: &models.Metrics{},
			error:  models.ErrUnknownMetricType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockIRepository(ctrl)
			s := &domain{
				repo: store,
			}

			if tt.prepare != nil {
				tt.prepare(store, tt.metric)
			}

			err := s.SetMetric(nil, tt.metric)
			assert.Equal(t, tt.error, err)
		})
	}
}
