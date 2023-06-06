package model_test

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsGetPreparedValue(t *testing.T) {
	tests := []struct {
		name   string
		fValue float64
		iValue int64
		m      *model.Metrics
		want   string
	}{
		{
			name:   model.MetricTypeGauge,
			fValue: 340255.4088704579,
			m: &model.Metrics{
				ID:    "test0",
				MType: model.MetricTypeGauge,
			},
			want: "340255.4088704579",
		},
		{
			name:   model.MetricTypeCounter,
			iValue: 42,
			m: &model.Metrics{
				ID:    "test1",
				MType: model.MetricTypeCounter,
			},
			want: "42",
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			if test.fValue != 0 {
				test.m.Value = &test.fValue
			}
			if test.iValue != 0 {
				test.m.Delta = &test.iValue
			}
			got, err := test.m.GetPreparedValue()
			assert.NoError(t, err)
			assert.Equalf(t, test.want, got, "problem test name: \"%s\"", test.name)
		})
	}
}

func TestMetricsUpdateByMetricUnit(t *testing.T) {
	tests := []struct {
		name       string
		fValueWant float64
		iValueWant int64
		mu         model.MetricUnit
		want       *model.Metrics
	}{
		{
			name: model.MetricTypeGauge,
			mu: model.MetricUnit{
				Type:       model.MetricTypeGauge,
				Name:       "test0",
				Value:      "3342.55",
				ValueFloat: 3342.55,
			},
			fValueWant: 3342.55,
			want: &model.Metrics{
				ID:    "test0",
				MType: model.MetricTypeGauge,
			},
		},
		{
			name: model.MetricTypeCounter,
			mu: model.MetricUnit{
				Type:       model.MetricTypeCounter,
				Name:       "test1",
				Value:      "42",
				ValueInt:   42,
				ValueFloat: 0,
			},
			iValueWant: 42,
			want: &model.Metrics{
				ID:    "test1",
				MType: model.MetricTypeCounter,
			},
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			if test.fValueWant != 0 {
				test.want.Value = &test.fValueWant
			}
			if test.iValueWant != 0 {
				test.want.Delta = &test.iValueWant
			}
			got := &model.Metrics{}
			got.UpdateByMetricUnit(test.mu)
			assert.Equalf(t, test.want, got, "problem test name: \"%s\"", test.name)
		})
	}
}
