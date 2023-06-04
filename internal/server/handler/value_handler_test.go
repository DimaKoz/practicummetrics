package handler

import (
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
			name:   "test 404 - 0",
			method: http.MethodGet,
			target: "/status",
			want: want{
				code:        http.StatusNotFound,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "test 404 - 1",
			method: http.MethodGet,
			target: "/value/gauge/testCounter132",
			want: want{
				code:        http.StatusNotFound,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "test 404 - 2",
			method: http.MethodGet,
			target: "/value/gauge/testCounter132",
			want: want{
				code:        http.StatusNotFound,
				response:    ``,
				contentType: "",
			},
		},
	}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(test.method, test.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			if i != len(tests)-1 {
				c.SetParamNames([]string{"name"}...)
				c.SetParamValues([]string{"testCounter132"}...)

			}
			_ = ValueHandler(c)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode, "StatusCode got: %v, want: %v", res.StatusCode, test.want.code)

			_ = res.Body.Close()

		})
	}
}
