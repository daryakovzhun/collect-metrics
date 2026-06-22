package handler

import (
	"fmt"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	typePath  = "metric_type"
	namePath  = "metric_name"
	valuePath = "metric_value"
)

func (h *Handler) SetMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := getMetricFromReq(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.domain.SetMetric(r.Context(), metric)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := getMetricFromReq(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err = h.domain.GetMetric(r.Context(), metric)
	if err != nil {
		handleError(w, err)
		return
	}

	var body []byte
	switch metric.MType {
	case models.Gauge:
		body = []byte(fmt.Sprintf("%g", *metric.Value))
	case models.Counter:
		body = []byte(fmt.Sprintf("%d", *metric.Delta))
	default:
		handleError(w, models.ErrUnknownMetricType)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(body); err != nil {
		handleError(w, err)
	}
}

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {

}

func getMetricFromReq(r *http.Request, parseValue bool) (*models.Metrics, error) {
	metric := models.Metrics{
		ID:    chi.URLParam(r, namePath),
		MType: chi.URLParam(r, typePath),
		Delta: nil, // counter
		Value: nil, // gauge
	}

	if len(metric.ID) == 0 {
		return nil, models.ErrEmptyMetricName
	}

	if parseValue {
		switch metric.MType {
		case models.Gauge:
			value, err := strconv.ParseFloat(chi.URLParam(r, valuePath), 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse metric gauge value: %w", err)
			}

			metric.Value = &value
		case models.Counter:
			delta, err := strconv.ParseInt(chi.URLParam(r, valuePath), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse metric counter value: %w", err)
			}

			metric.Delta = &delta
		default:
			return nil, models.ErrUnknownMetricType
		}
	}

	return &metric, nil
}
