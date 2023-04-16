package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "test counter StatusOk",
			request: "/update/counter/testCounter/100",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:    "test counter no value",
			request: "/update/counter/testCounter/none",
			want: want{
				code:        http.StatusBadRequest,
				response:    ``,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, test.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			paramValues := strings.Split(test.request, "/")

			c.SetParamNames([]string{"type", "name", "value"}...)
			c.SetParamValues(paramValues[2:]...)

			_ = UpdateHandler(c)

			res := w.Result()
			// проверяем код ответа
			got := res.StatusCode
			assert.Equal(t, test.want.code, got, "StatusCode got: %v, want: %v", got, test.want.code)

			_ = res.Body.Close()

		})
	}
}

/*func Test_processPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  model.MetricUnit
		want1 *error2.RequestError
	}{
		{name: "empty path",
			args: args{
				path: "",
			},
			want:  model.MetricUnit{},
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")},
		},
		{name: "wrong number '/'",
			args: args{
				path: "update/zs/",
			},
			want:  model.MetricUnit{},
			want1: &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")},
		},
		{name: "wrong value",
			args: args{
				path: "/update/counter/me/none",
			},
			want:  model.MetricUnit{},
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("bad value")},
		},
		{name: "ok counter metric",
			args: args{
				path: "/update/counter/me/42",
			},
			want:  model.MetricUnit{Type: "counter", Name: "me", Value: "42", ValueInt: 42, ValueFloat: 0},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := processPath(tt.args.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processPath() got = %v, want %v", got, tt.want)
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
*/
