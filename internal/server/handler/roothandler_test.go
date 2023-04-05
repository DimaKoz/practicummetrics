package handler

import (
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
			request := httptest.NewRequest(http.MethodGet, "/status", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			RootHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			if res.StatusCode != test.want.code {
				t.Errorf("StatusCode got: %v, want: %v", res.StatusCode, test.want.code)
			}

			res.Body.Close()

		})
	}
}
