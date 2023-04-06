package error

import (
	"errors"
	"net/http"
	"testing"
)

func TestRequestError_Error(t *testing.T) {
	type fields struct {
		StatusCode int
		Err        error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "test error()",
			fields: fields{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")},
			want:   "unavailable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := &RequestError{
				StatusCode: tt.fields.StatusCode,
				Err:        tt.fields.Err,
			}
			if got := se.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCode(t *testing.T) {
	type args struct {
		se RequestError
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "test status code",
			args: args{RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCode(tt.args.se); got != tt.want {
				t.Errorf("StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
