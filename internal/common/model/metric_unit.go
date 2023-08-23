package model

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// MetricTypeGauge represents "gauge" MetricUnit.Type.
	MetricTypeGauge = "gauge"
	// MetricTypeCounter represents "counter" MetricUnit.Type.
	MetricTypeCounter = "counter"
)

// EmptyMetric is an empty struct of MetricUnit.
var EmptyMetric = MetricUnit{} //nolint:exhaustruct

// MetricUnit represents a metric.
type MetricUnit struct {
	Type       string
	Name       string
	Value      string
	ValueInt   int64
	ValueFloat float64
}

// ErrUnknownType represents an error with an unknown type of the metric.
var ErrUnknownType = errors.New("unknown metric type") // should use for StatusCode: http.StatusNotImplemented

// ErrEmptyValue represents an error which related to empty MetricUnit.Name and/or MetricUnit.Value.
// fo StatusCode: http.StatusBadRequest.
var ErrEmptyValue = errors.New("to create a metric you must provide `name` and `value`")

// NewMetricUnit creates an instance of MetricUnit or returns an error.
func NewMetricUnit(metricType string, metricName string, metricValue string) (MetricUnit, error) {
	if metricType != MetricTypeGauge && metricType != MetricTypeCounter {
		return EmptyMetric, ErrUnknownType
	}

	if metricName == "" || metricValue == "" {
		return EmptyMetric, ErrEmptyValue
	}
	result := MetricUnit{} //nolint:exhaustruct
	result.Type = metricType
	result.Name = metricName
	result.Value = metricValue

	if metricType == MetricTypeGauge {
		if s, err := strconv.ParseFloat(metricValue, 64); err == nil {
			result.ValueFloat = s
		} else {
			err = fmt.Errorf("bad value: failed to parse metricValue by: %w", err) // StatusCode: http.StatusBadRequest

			return EmptyMetric, err
		}
	}

	if metricType == MetricTypeCounter {
		if s, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
			result.ValueInt = s
		} else {
			err = fmt.Errorf("bad value: failed to parse metricValue by: %w", err) // StatusCode: http.StatusBadRequest

			return EmptyMetric, err
		}
	}

	return result, nil
}

// Clone creates a copy of instance of MetricUnit and returns it.
func (mu MetricUnit) Clone() MetricUnit {
	return MetricUnit{
		Type:       mu.Type,
		Name:       mu.Name,
		Value:      mu.Value,
		ValueFloat: mu.ValueFloat,
		ValueInt:   mu.ValueInt,
	}
}

// GetPath gets a part of path by MetricUnit and returns it.
func (mu MetricUnit) GetPath() string {
	return fmt.Sprintf("%s/%s/%s", mu.Type, mu.Name, mu.Value)
}
