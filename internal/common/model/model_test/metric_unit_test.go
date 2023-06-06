package model_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
)

func TestNewMetricUnit(t *testing.T) {
	errBadValue := errors.New("bad value")

	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name    string
		args    args
		want    model.MetricUnit
		wantErr error
	}{
		{
			name: "normal counter",
			args: args{
				metricType:  model.MetricTypeCounter,
				metricName:  "test",
				metricValue: "42",
			},
			want:    model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "42", ValueInt: 42, ValueFloat: 0},
			wantErr: nil,
		},
		{
			name: "normal gauge",
			args: args{
				metricType:  model.MetricTypeGauge,
				metricName:  "test",
				metricValue: "42",
			},
			want:    model.MetricUnit{Type: model.MetricTypeGauge, Name: "test", Value: "42", ValueInt: 0, ValueFloat: 42},
			wantErr: nil,
		},
		{
			name: "unknown type",
			args: args{
				metricType:  "xyz",
				metricName:  "test",
				metricValue: "42",
			},
			want:    model.EmptyMetric,
			wantErr: model.ErrorUnknownType,
		},
		{
			name: "empty name",
			args: args{
				metricType:  model.MetricTypeGauge,
				metricName:  "",
				metricValue: "42",
			},
			want:    model.EmptyMetric,
			wantErr: model.ErrorEmptyValue,
		},
		{
			name: "empty value",
			args: args{
				metricType:  model.MetricTypeGauge,
				metricName:  "qaz",
				metricValue: "",
			},
			want:    model.EmptyMetric,
			wantErr: model.ErrorEmptyValue,
		},
		{
			name: "no float value",
			args: args{
				metricType:  model.MetricTypeGauge,
				metricName:  "qaz",
				metricValue: "xexe",
			},
			want:    model.EmptyMetric,
			wantErr: errBadValue,
		},
		{
			name: "no int value",
			args: args{
				metricType:  model.MetricTypeCounter,
				metricName:  "qaz",
				metricValue: "xexe",
			},
			want:    model.EmptyMetric,
			wantErr: errBadValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := model.NewMetricUnit(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricUnit() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != got1 && !strings.Contains(tt.wantErr.Error(), "bad value") {
				t.Errorf("processPath() got1 = %v, want %v", got1, tt.wantErr)
			}
		})
	}
}

func TestMetricUnitClone(t *testing.T) {
	tests := []struct {
		name string
		pass model.MetricUnit
		want model.MetricUnit
	}{
		{
			name: "clone",
			pass: model.MetricUnit{
				Type:       model.MetricTypeGauge,
				Name:       "heap",
				ValueInt:   0,
				ValueFloat: 4932.99,
				Value:      "4932.99",
			},
			want: model.MetricUnit{
				Type:       model.MetricTypeGauge,
				Name:       "heap",
				ValueInt:   0,
				ValueFloat: 4932.99,
				Value:      "4932.99",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pass.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricUnitGetPath(t *testing.T) {
	tests := []struct {
		name string
		mu   model.MetricUnit
		want string
	}{
		{
			name: "a path from MetricUnit",
			mu: model.MetricUnit{
				Type:       model.MetricTypeGauge,
				Name:       "b",
				Value:      "42",
				ValueInt:   0,
				ValueFloat: 42,
			},
			want: "gauge/b/42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mu.GetPath(); got != tt.want {
				t.Errorf("GetPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
