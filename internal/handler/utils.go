package handler

import (
	"errors"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"go.uber.org/zap"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, models.ErrUnknownMetricType):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	logger.Log.Error("handle error", zap.Error(err))
}

func validateMetric(metric models.Metrics) error {
	if len(metric.ID) == 0 {
		return models.ErrEmptyMetricName
	}

	switch metric.MType {
	case models.Gauge:
		if metric.Value == nil {
			return models.ErrValueNotSet
		}
	case models.Counter:
		if metric.Delta == nil {
			return models.ErrValueNotSet
		}
	default:
		return models.ErrUnknownMetricType
	}

	return nil
}
