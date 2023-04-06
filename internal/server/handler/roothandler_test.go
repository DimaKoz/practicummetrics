package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test 404",
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
			request := httptest.NewRequest(http.MethodGet, "/status", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			RootHandler(c)

			res := w.Result()
			// проверяем код ответа
			if res.StatusCode != test.want.code {
				t.Errorf("StatusCode got: %v, want: %v", res.StatusCode, test.want.code)
			}

			res.Body.Close()

		})
	}
}
