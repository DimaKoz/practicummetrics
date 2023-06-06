package handler_test

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestValueHandlerJSON(t *testing.T) {
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
			request: "/value",
			reqJSON: "",
			want: want{
				code:        http.StatusNotFound,
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
			name:    "test OK",
			request: "/update",
			reqJSON: `{"id":"GetSet187"}`,
			want: want{
				code:        http.StatusOK,
				response:    "{\"id\":\"GetSet187\",\"type\":\"gauge\",\"value\":42}\n",
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			fileStorage := filepath.Join(tmpDir, "rep_values.txt")
			repository.SetupFilePathStorage(fileStorage)
			assert.NoError(t, repository.Save())
			mu, err := model.NewMetricUnit("gauge", "GetSet187", "42")
			assert.NoError(t, err)
			repository.AddMetric(mu)
			t.Cleanup(func() {
				_ = repository.Load()
			})

			echoFramework := echo.New()
			var body io.Reader
			if test.reqJSON != "" {
				body = strings.NewReader(test.reqJSON)
			}

			request := httptest.NewRequest(http.MethodPost, test.request, body)

			request.Header.Set("Content-Type", "application/json")

			respRecorder := httptest.NewRecorder()
			c := echoFramework.NewContext(request, respRecorder)

			err = handler.ValueHandlerJSON(c)
			assert.NoError(t, err, "expected no errors")

			res := respRecorder.Result()
			defer res.Body.Close()
			// проверяем код ответа
			got := res.StatusCode

			assert.Equal(t, test.want.code, got, "StatusCode got: %v, want: %v", got, test.want.code)
			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "expected no errors")
			if got == http.StatusOK {
				assert.Equal(t, test.want.response, string(b))
			}
		})
	}
}
