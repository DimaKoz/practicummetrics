package model

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	// MetricTypeGauge represents "gauge" MetricUnit.Type
	MetricTypeGauge = "gauge"
	// MetricTypeCounter represents "counter" MetricUnit.Type
	MetricTypeCounter = "counter"
)

var EmptyMetric = MetricUnit{}

// MetricUnit represents a metric
type MetricUnit struct {
	Type       string
	Name       string
	Value      string
	ValueInt   int64
	ValueFloat float64
}

// ErrorUnknownType represents an error with an unknown type of the metric
var ErrorUnknownType = errors.New("unknown metric type") //should use for StatusCode: http.StatusNotImplemented

// ErrorEmptyValue represents an error which related to empty MetricUnit.Name and/or MetricUnit.Value
var ErrorEmptyValue = errors.New("to create a metric you must provide `name` and `value`") // StatusCode: http.StatusBadRequest

// NewMetricUnit creates an instance of MetricUnit or returns an error
func NewMetricUnit(metricType string, metricName string, metricValue string) (MetricUnit, error) {
	if metricType != MetricTypeGauge && metricType != MetricTypeCounter {
		return EmptyMetric, ErrorUnknownType
	}
	if metricName == "" || metricValue == "" {
		return EmptyMetric, ErrorEmptyValue
	}
	var result = MetricUnit{}
	result.Type = metricType
	result.Name = metricName
	result.Value = metricValue

	if metricType == MetricTypeGauge {
		if s, err := strconv.ParseFloat(metricValue, 64); err == nil {
			result.ValueFloat = s
		} else {
			return EmptyMetric, fmt.Errorf("bad value: failed to parse metricValue by: %w", err) // StatusCode: http.StatusBadRequest
		}
	}

	if metricType == MetricTypeCounter {
		if s, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
			result.ValueInt = s
		} else {
			return EmptyMetric, fmt.Errorf("bad value: failed to parse metricValue by: %w", err) // StatusCode: http.StatusBadRequest
		}
	}
	return result, nil

}

func (mu MetricUnit) Clone() MetricUnit {
	return MetricUnit{
		Type:       mu.Type,
		Name:       mu.Name,
		Value:      mu.Value,
		ValueFloat: mu.ValueFloat,
		ValueInt:   mu.ValueInt,
	}
}

func (mu MetricUnit) GetPath() string {
	return fmt.Sprintf("%s/%s/%s", mu.Type, mu.Name, mu.Value)
}
