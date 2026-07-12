package models

import "errors"

var (
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrEmptyMetricName   = errors.New("empty metric name")
	ErrNotFound          = errors.New("recourse not found")
)
