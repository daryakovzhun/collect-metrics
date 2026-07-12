package handler

import (
	"errors"
	"fmt"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
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

	fmt.Fprintln(w, err.Error())
}
