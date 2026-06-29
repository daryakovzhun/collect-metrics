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

	err = h.domain.SetMetric(r.Context(), &metric)
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

	metric, err = h.domain.GetMetric(r.Context(), &metric)
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
	metrics, err := h.domain.GetAllMetrics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	views := make([]MetricView, 0, len(metrics))
	for _, m := range metrics {
		view := MetricView{
			Name: m.ID,
			Type: m.MType,
		}
		// Формируем строковое значение
		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				view.Value = fmt.Sprintf("%g", *m.Value)
			} else {
				view.Value = "null"
			}
		case models.Counter:
			if m.Delta != nil {
				view.Value = fmt.Sprintf("%d", *m.Delta)
			} else {
				view.Value = "null"
			}
		default:
			view.Value = "unknown type"
		}
		views = append(views, view)
	}

	data := MetricsList{Items: views}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = h.metricsTmpl.Execute(w, data)
	if err != nil {
		handleError(w, err)
		return
	}
}

func getMetricFromReq(r *http.Request, parseValue bool) (models.Metrics, error) {
	metric := models.Metrics{
		ID:    chi.URLParam(r, namePath),
		MType: chi.URLParam(r, typePath),
		Delta: nil, // counter
		Value: nil, // gauge
	}

	if len(metric.ID) == 0 {
		return models.Metrics{}, models.ErrEmptyMetricName
	}

	if parseValue {
		switch metric.MType {
		case models.Gauge:
			value, err := strconv.ParseFloat(chi.URLParam(r, valuePath), 64)
			if err != nil {
				return models.Metrics{}, fmt.Errorf("failed to parse metric gauge value: %w", err)
			}

			metric.Value = &value
		case models.Counter:
			delta, err := strconv.ParseInt(chi.URLParam(r, valuePath), 10, 64)
			if err != nil {
				return models.Metrics{}, fmt.Errorf("failed to parse metric counter value: %w", err)
			}

			metric.Delta = &delta
		default:
			return models.Metrics{}, models.ErrUnknownMetricType
		}
	}

	return metric, nil
}
