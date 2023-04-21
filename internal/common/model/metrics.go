package model

import "strconv"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) UpdateByMetricUnit(mu MetricUnit) {
	m.ID = mu.Name
	m.MType = mu.Type
	if mu.Type == MetricTypeGauge {
		m.Value = &mu.ValueFloat
	} else {
		m.Delta = &mu.ValueInt
	}
}

func (m *Metrics) GetPreparedValue() string {
	var metricValue string
	if m.MType == MetricTypeGauge {
		metricValue = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	} else {
		metricValue = strconv.FormatInt(*m.Delta, 10)
	}
	return metricValue
}
