package handler

import (
	"fmt"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"net/http"
	"strconv"
)

const (
	typePath  = "metric_type"
	namePath  = "metric_name"
	valuePath = "metric_value"
)

func (h *Handler) SetMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := getMetricFromReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.domain.SetMetric(r.Context(), metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getMetricFromReq(r *http.Request) (*models.Metrics, error) {
	metric := models.Metrics{
		ID:    r.PathValue(namePath),
		MType: r.PathValue(typePath),
		Delta: nil, // counter
		Value: nil, // gauge
	}

	if len(metric.ID) == 0 {
		return nil, models.ErrEmptyMetricName
	}

	switch metric.MType {
	case models.Gauge:
		value, err := strconv.ParseFloat(r.PathValue(valuePath), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse metric gauge value: %w", err)
		}

		metric.Value = &value
	case models.Counter:
		delta, err := strconv.ParseInt(r.PathValue(valuePath), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse metric counter value: %w", err)
		}

		metric.Delta = &delta
	default:
		return nil, models.ErrUnknownMetricType
	}

	return &metric, nil
}
