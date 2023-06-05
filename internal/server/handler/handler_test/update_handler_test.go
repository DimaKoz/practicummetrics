package handler_test

import (
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
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
			c.SetPath("/update/:type/:name/:value")
			c.SetParamNames([]string{"type", "name", "value"}...)
			c.SetParamValues(paramValues[2:]...)

			assert.NoError(t, handler.UpdateHandler(c), "expected no errors")

			res := w.Result()
			// проверяем код ответа
			got := res.StatusCode
			assert.Equal(t, test.want.code, got, "StatusCode got: %v, want: %v", got, test.want.code)

			_ = res.Body.Close()

		})
	}
}
