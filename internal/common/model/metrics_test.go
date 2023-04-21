package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetrics_GetPreparedValue(t *testing.T) {
	tests := []struct {
		name   string
		fValue float64
		iValue int64
		m      *Metrics
		want   string
	}{
		{
			name:   MetricTypeGauge,
			fValue: 340255.4088704579,
			m: &Metrics{
				ID:    "test0",
				MType: MetricTypeGauge,
			},
			want: "340255.4088704579",
		},
		{
			name:   MetricTypeCounter,
			iValue: 42,
			m: &Metrics{
				ID:    "test1",
				MType: MetricTypeCounter,
			},
			want: "42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fValue != 0 {
				tt.m.Value = &tt.fValue
			}
			if tt.iValue != 0 {
				tt.m.Delta = &tt.iValue
			}
			got := tt.m.GetPreparedValue()
			assert.Equalf(t, tt.want, got, "problem test name: \"%s\"", tt.name)
		})
	}
}

func TestMetrics_UpdateByMetricUnit(t *testing.T) {
	tests := []struct {
		name       string
		fValueWant float64
		iValueWant int64
		mu         MetricUnit
		want       *Metrics
	}{
		{
			name: MetricTypeGauge,
			mu: MetricUnit{
				Type:       MetricTypeGauge,
				Name:       "test0",
				Value:      "3342.55",
				ValueFloat: 3342.55,
			},
			fValueWant: 3342.55,
			want: &Metrics{
				ID:    "test0",
				MType: MetricTypeGauge,
			},
		},
		{
			name: MetricTypeCounter,
			mu: MetricUnit{
				Type:     MetricTypeCounter,
				Name:     "test1",
				Value:    "42",
				ValueInt: 42,
			},
			iValueWant: 42,
			want: &Metrics{
				ID:    "test1",
				MType: MetricTypeCounter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fValueWant != 0 {
				tt.want.Value = &tt.fValueWant
			}
			if tt.iValueWant != 0 {
				tt.want.Delta = &tt.iValueWant
			}
			got := &Metrics{}
			got.UpdateByMetricUnit(tt.mu)
			assert.Equalf(t, tt.want, got, "problem test name: \"%s\"", tt.name)
		})
	}
}
