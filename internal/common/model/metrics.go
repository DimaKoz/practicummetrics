package model

import (
	"errors"
	"strconv"
)

// Metrics represents a metric from JSON.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewEmptyMetrics() *Metrics {
	result := &Metrics{ID: "", MType: "", Delta: nil, Value: nil}

	return result
}

func (m *Metrics) UpdateByMetricUnit(metricUnit MetricUnit) {
	m.ID = metricUnit.Name
	m.MType = metricUnit.Type

	if metricUnit.Type == MetricTypeGauge {
		m.Value = &metricUnit.ValueFloat
		m.Delta = nil
	} else {
		m.Delta = &metricUnit.ValueInt
		m.Value = nil
	}
}

var (
	errProblemDelta = errors.New("couldn't convert Metrics.Delta to a string, it must not be nil")
	errProblemValue = errors.New("couldn't convert Metrics.Value to a string, it must not be nil")
)

func (m *Metrics) GetPreparedValue() (string, error) {
	var metricValue string

	if m.MType == MetricTypeGauge {
		if m.Value == nil {
			return "", errProblemValue
		}
		metricValue = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	} else {
		if m.Delta == nil {
			return "", errProblemDelta
		}
		metricValue = strconv.FormatInt(*m.Delta, 10)
	}

	return metricValue, nil
}
