package model

type Metrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) Convert(mu MetricUnit) {
	m.ID = mu.Name
	m.MType = mu.Type
	if mu.Type == MetricTypeGauge {
		m.Value = mu.ValueFloat
	} else {
		m.Delta = mu.ValueInt
	}
}
