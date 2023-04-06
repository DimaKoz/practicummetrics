package repository

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
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
		wantkey string
		want    *model.MetricUnit
	}{
		{name: "counter",
			args: []args{
				{mu: model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "42", ValueI: 42, ValueF: 0}},
				{mu: model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "10", ValueI: 10, ValueF: 0}},
			},
			wantkey: "test",
			want:    &model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "52", ValueI: 52, ValueF: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, unit := range tt.args {
				AddMetricMemStorage(unit.mu)
			}
			if got, ok := instanceMemSt.storage[tt.wantkey]; ok {
				if !reflect.DeepEqual(&got, tt.want) {
					t.Errorf("AddMetricMemStorage() got = %v, want %v", got, tt.want)
				}
			} else {
				t.Errorf("not found stored result")
			}
		})
	}

}
