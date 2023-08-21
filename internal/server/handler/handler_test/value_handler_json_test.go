package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
			name: "test no json", request: "/value",
			reqJSON: "", want: want{code: http.StatusNotFound, response: ``, contentType: ""},
		},
		{
			name: "test bad json", request: "/update",
			reqJSON: "{", want: want{code: http.StatusBadRequest, response: ``, contentType: ""},
		},
		{
			name: "test OK", request: "/update", reqJSON: `{"id":"GetSet187"}`,
			want: want{
				code:        http.StatusOK,
				response:    "{\"id\":\"GetSet187\",\"type\":\"gauge\",\"value\":42}\n",
				contentType: "",
			},
		},
	}

	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			fileStorage := filepath.Join(tmpDir, "rep_values.txt")
			repository.SetupFilePathStorage(fileStorage)
			assert.NoError(t, repository.SaveVariant())
			mu, err := model.NewMetricUnit("gauge", "GetSet187", "42")
			assert.NoError(t, err)
			repository.AddMetric(mu)
			t.Cleanup(func() { _ = repository.Load() })
			echoFramework := echo.New()
			request, respRecorder := setupEchoHandlerJSONTest(test.reqJSON, test.request)
			ctx := echoFramework.NewContext(request, respRecorder)
			err = handler.NewBaseHandler(nil).ValueHandlerJSON(ctx)
			assert.NoError(t, err, "expected no errors")

			res := respRecorder.Result()
			defer res.Body.Close()
			got := res.StatusCode // проверяем код ответа
			assert.Equal(t, test.want.code, got, "StatusCode got: %v, want: %v", got, test.want.code)
			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "expected no errors")
			if got == http.StatusOK {
				assert.Equal(t, test.want.response, string(b))
			}
		})
	}
}

func setupEchoHandlerJSONTest(reqJSON string, req string) (*http.Request, *httptest.ResponseRecorder) {
	var body io.Reader
	if reqJSON != "" {
		body = strings.NewReader(reqJSON)
	}

	request := httptest.NewRequest(http.MethodPost, req, body)

	request.Header.Set("Content-Type", "application/json")

	respRecorder := httptest.NewRecorder()

	return request, respRecorder
}
