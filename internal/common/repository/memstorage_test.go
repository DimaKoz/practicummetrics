package repository

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAddMetricMemStorage(t *testing.T) {

	type args struct {
		mu model.MetricUnit
	}
	tests := []struct {
		name    string
		args    []args
		wantKey string
		want    *model.MetricUnit
	}{
		{name: "counter",
			args: []args{
				{mu: model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "42", ValueInt: 42, ValueFloat: 0}},
				{mu: model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "10", ValueInt: 10, ValueFloat: 0}},
			},
			wantKey: "test",
			want:    &model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "52", ValueInt: 52, ValueFloat: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, unit := range tt.args {
				AddMetric(unit.mu)
			}
			if got, ok := memStorage.storage[tt.wantKey]; ok {
				if !reflect.DeepEqual(&got, tt.want) {
					t.Errorf("AddMetric() got = %v, want %v", got, tt.want)
				}
			} else {
				t.Errorf("not found stored result")
			}
		})
	}

}

func TestGetMetricByName(t *testing.T) {
	type args struct {
		search string
		add    []model.MetricUnit
	}
	tests := []struct {
		name string
		args args
		want *model.MetricUnit
	}{
		{
			name: "empty key",
			args: args{
				search: "",
				add:    []model.MetricUnit{},
			},
			want: nil,
		},
		{
			name: "wanted key",
			args: args{
				search: "wanted",
				add: []model.MetricUnit{
					{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
					{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
				},
			},
			want: &model.MetricUnit{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := memStorage.storage
			memStorage.storage = make(map[string]model.MetricUnit, 0)
			t.Cleanup(func() { memStorage.storage = orig })
			for _, v := range tt.args.add {
				AddMetric(v)
			}
			got := GetMetricByName(tt.args.search)
			assert.Equal(t, got, tt.want, "GetMetricByName() = %v, want %v", got, tt.want)
		})
	}
}

func TestGetMetricsMemStorage(t *testing.T) {

	tests := []struct {
		add  []model.MetricUnit
		name string
		want []model.MetricUnit
	}{
		{
			name: "empty",
			add:  []model.MetricUnit{},
			want: []model.MetricUnit{},
		}, {name: "wanted key",
			add: []model.MetricUnit{
				{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
				{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
			},
			want: []model.MetricUnit{
				{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
				{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := memStorage.storage
			memStorage.storage = make(map[string]model.MetricUnit, 0)
			t.Cleanup(func() { memStorage.storage = orig })
			for _, v := range tt.add {
				AddMetric(v)
			}
			assert.ElementsMatch(t, tt.want, GetAllMetrics(), "GetAllMetrics()")
		})
	}
}
