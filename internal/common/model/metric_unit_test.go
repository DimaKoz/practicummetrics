package model

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestNewMetricUnit(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name    string
		args    args
		want    MetricUnit
		wantErr error
	}{
		{name: "normal counter",
			args: args{
				metricType:  MetricTypeCounter,
				metricName:  "test",
				metricValue: "42",
			},
			want:    MetricUnit{MetricTypeCounter, "test", "42", 42, 0},
			wantErr: nil,
		},
		{name: "normal gauge",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "test",
				metricValue: "42",
			},
			want:    MetricUnit{MetricTypeGauge, "test", "42", 0, 42},
			wantErr: nil,
		},
		{name: "unknown type",
			args: args{
				metricType:  "xyz",
				metricName:  "test",
				metricValue: "42",
			},
			want:    EmptyMetric,
			wantErr: ErrorUnknownType,
		},
		{name: "empty name",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "",
				metricValue: "42",
			},
			want:    EmptyMetric,
			wantErr: ErrorEmptyValue,
		},
		{name: "empty value",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "qaz",
				metricValue: "",
			},
			want:    EmptyMetric,
			wantErr: ErrorEmptyValue,
		},
		{name: "no float value",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "qaz",
				metricValue: "xexe",
			},
			want:    EmptyMetric,
			wantErr: errors.New("bad value"),
		},
		{name: "no int value",
			args: args{
				metricType:  MetricTypeCounter,
				metricName:  "qaz",
				metricValue: "xexe",
			},
			want:    EmptyMetric,
			wantErr: errors.New("bad value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewMetricUnit(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
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
		pass MetricUnit
		want MetricUnit
	}{
		{
			name: "clone",
			pass: MetricUnit{
				Type:       MetricTypeGauge,
				Name:       "heap",
				ValueInt:   0,
				ValueFloat: 4932.99,
				Value:      "4932.99",
			},
			want: MetricUnit{
				Type:       MetricTypeGauge,
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
		mu   MetricUnit
		want string
	}{
		{
			name: "a path from MetricUnit",
			mu: MetricUnit{
				Type:       MetricTypeGauge,
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
