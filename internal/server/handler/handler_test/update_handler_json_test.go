package handler_test

import (
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateHandlerJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		reqJSON string
		want    want
	}{

		{
			name:    "test no json",
			request: "/update",
			reqJSON: "",
			want: want{
				code:        http.StatusBadRequest,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:    "test bad json",
			request: "/update",
			reqJSON: "{",
			want: want{
				code:        http.StatusBadRequest,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:    "test wrong type metrics",
			request: "/update",
			reqJSON: "{\"id\":\"GetSet185\",\"type\":\"gauge1\",\"delta\":1}",
			want: want{
				code:        http.StatusNotImplemented,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:    "test OK",
			request: "/update",
			reqJSON: `{"id":"GetSet186","type":"counter","delta":1}`,
			want: want{
				code:        http.StatusOK,
				response:    "{\"id\":\"GetSet186\",\"type\":\"counter\",\"delta\":1}\n",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			var body io.Reader = nil
			if tt.reqJSON != "" {
				body = strings.NewReader(tt.reqJSON)
			}

			request := httptest.NewRequest(http.MethodPost, tt.request, body)

			request.Header.Set("Content-Type", "application/json")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)

			err := handler.UpdateHandlerJSON(c)
			assert.NoError(t, err, "expected no errors")

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			got := res.StatusCode

			assert.Equal(t, tt.want.code, got, "StatusCode got: %v, want: %v", got, tt.want.code)
			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "expected no errors")
			if got == http.StatusOK {
				assert.Equal(t, tt.want.response, string(b))
			}

		})

	}
}
