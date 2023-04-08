package handler

import (
	"errors"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValueHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		target string
		want   want
	}{
		{
			name:   "test 404",
			method: http.MethodGet,
			target: "/status",
			want: want{
				code:        http.StatusNotFound,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "test 404",
			method: http.MethodGet,
			target: "/value/gauge/testCounter132",
			want: want{
				code:        http.StatusNotFound,
				response:    ``,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(test.method, test.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			ValueHandler(c)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode, "StatusCode got: %v, want: %v", res.StatusCode, test.want.code)

			res.Body.Close()

		})
	}
}

func Test_getNameFromPath(t *testing.T) {

	tests := []struct {
		name  string
		path  string
		want  string
		want1 *error2.RequestError
	}{
		{
			name:  "empty path",
			path:  "",
			want:  "",
			want1: &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")},
		},
		{
			name:  "bad path",
			path:  "/6/5/4/3/2/1",
			want:  "",
			want1: &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")},
		},
		{
			name:  "ok path",
			path:  "/value/gauge/testCounter",
			want:  "testCounter",
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getNameFromPath(tt.path)
			assert.Equalf(t, tt.want, got, "getNameFromPath(%v)", tt.path)
			assert.Equalf(t, tt.want1, got1, "getNameFromPath(%v)", tt.path)
		})
	}
}
