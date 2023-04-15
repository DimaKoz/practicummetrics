package model

import (
	"errors"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"net/http"
	"reflect"
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
		name  string
		args  args
		want  MetricUnit
		want1 *error2.RequestError
	}{
		{name: "normal counter",
			args: args{
				metricType:  MetricTypeCounter,
				metricName:  "test",
				metricValue: "42",
			},
			want:  MetricUnit{MetricTypeCounter, "test", "42", 42, 0},
			want1: nil,
		},
		{name: "normal gauge",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "test",
				metricValue: "42",
			},
			want:  MetricUnit{MetricTypeGauge, "test", "42", 0, 42},
			want1: nil,
		},
		{name: "unknown type",
			args: args{
				metricType:  "xyz",
				metricName:  "test",
				metricValue: "42",
			},
			want:  EmptyMetric,
			want1: &error2.RequestError{StatusCode: http.StatusNotImplemented, Err: errors.New("unknown type")},
		},
		{name: "empty name",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "",
				metricValue: "42",
			},
			want:  EmptyMetric,
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")},
		},
		{name: "empty value",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "qaz",
				metricValue: "",
			},
			want:  EmptyMetric,
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")},
		},
		{name: "no float value",
			args: args{
				metricType:  MetricTypeGauge,
				metricName:  "qaz",
				metricValue: "xexe",
			},
			want:  EmptyMetric,
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("bad value")},
		},
		{name: "no int value",
			args: args{
				metricType:  MetricTypeCounter,
				metricName:  "qaz",
				metricValue: "xexe",
			},
			want:  EmptyMetric,
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("bad value")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewMetricUnit(tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricUnit() got = %v, want %v", got, tt.want)
			}
			if tt.want1 != nil && got1 != nil {
				if got1.StatusCode != tt.want1.StatusCode {
					t.Errorf("processPath() got1 = %v, want %v", got1, tt.want1)
				}
			} else if tt.want1 != got1 {
				t.Errorf("processPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMetricUnit_Clone(t *testing.T) {
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

func TestMetricUnit_GetPath(t *testing.T) {

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
