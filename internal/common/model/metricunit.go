package model

import (
	"errors"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"net/http"
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

// NewMetricUnit creates an instance of MetricUnit or returns *error2.RequestError
func NewMetricUnit(metricType string, metricName string, metricValue string) (MetricUnit, *error2.RequestError) {
	if metricType != MetricTypeGauge && metricType != MetricTypeCounter {
		return EmptyMetric, &error2.RequestError{StatusCode: http.StatusNotImplemented, Err: errors.New("unknown type")}
	}
	if metricName == "" || metricValue == "" {
		return EmptyMetric, &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}
	}
	var result = MetricUnit{}
	result.Type = metricType
	result.Name = metricName
	result.Value = metricValue

	if metricType == MetricTypeGauge {
		if s, err := strconv.ParseFloat(metricValue, 64); err == nil {
			result.ValueFloat = s
		} else {
			return EmptyMetric, &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("bad value")}
		}
	}

	if metricType == MetricTypeCounter {
		if s, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
			result.ValueInt = s
		} else {
			return EmptyMetric, &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("bad value")}
		}
	}
	var err *error2.RequestError = nil
	return result, err

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
	return mu.Type + "/" + mu.Name + "/" + mu.Value
}
