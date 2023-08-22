package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/DimaKoz/practicummetrics/internal/server/handler"
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
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, test.request, nil)
			responseRecorder := httptest.NewRecorder() // создаём новый Recorder
			ctx := e.NewContext(request, responseRecorder)
			paramValues := strings.Split(test.request, "/")
			ctx.SetPath("/update/:type/:name/:value")
			ctx.SetParamNames([]string{"type", "name", "value"}...)
			ctx.SetParamValues(paramValues[2:]...)

			assert.NoError(t, handler.UpdateHandler(ctx), "expected no errors")

			res := responseRecorder.Result()
			// проверяем код ответа
			got := res.StatusCode
			assert.Equal(t, test.want.code, got, "StatusCode got: %v, want: %v", got, test.want.code)

			_ = res.Body.Close()
		})
	}
}
