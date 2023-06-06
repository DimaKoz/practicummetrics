package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
			responseRecorder := httptest.NewRecorder() // создаём новый Recorder
			ctx := e.NewContext(request, responseRecorder)
			if i != len(tests)-1 {
				ctx.SetParamNames([]string{"name"}...)
				ctx.SetParamValues([]string{"testCounter132"}...)
			}
			_ = handler.ValueHandler(ctx)
			res := responseRecorder.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode, "StatusCode got: %v, want: %v", res.StatusCode, test.want.code)

			_ = res.Body.Close()
		})
	}
}
