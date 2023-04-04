package handler

import (
	"errors"
	error2 "github.com/DimaKoz/practicummetrics/internal/server/error"
	"github.com/DimaKoz/practicummetrics/internal/server/model"
	"net/http"
	"reflect"
	"testing"
)

func Test_processPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  *model.MetricUnit
		want1 *error2.RequestError
	}{
		{name: "empty path",
			args: args{
				path: "",
			},
			want:  nil,
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")},
		},
		{name: "wrong number '/'",
			args: args{
				path: "update/zs/",
			},
			want:  nil,
			want1: &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := processPath(tt.args.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processPath() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("processPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
