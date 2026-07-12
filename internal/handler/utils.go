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
