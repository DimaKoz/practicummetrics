package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type wantUpdHandlerJSON struct {
	code        int
	response    string
	contentType string
}

func TestUpdateHandlerJSON(t *testing.T) {
	tests := []struct {
		name    string
		request string
		reqJSON string
		want    wantUpdHandlerJSON
	}{
		{
			name: "test no json", request: "/update", reqJSON: "",
			want: wantUpdHandlerJSON{code: http.StatusBadRequest, response: ``, contentType: ""},
		},
		{
			name: "test bad json", request: "/update", reqJSON: "{",
			want: wantUpdHandlerJSON{code: http.StatusBadRequest, response: ``, contentType: ""},
		},
		{
			name: "test wrong type metrics", request: "/update",
			reqJSON: "{\"id\":\"GetSet185\",\"type\":\"gauge1\",\"delta\":1}",
			want:    wantUpdHandlerJSON{code: http.StatusNotImplemented, response: ``, contentType: ""},
		},
		{
			name: "test OK", request: "/update", reqJSON: `{"id":"GetSet186","type":"counter","delta":1}`,
			want: wantUpdHandlerJSON{
				code:        http.StatusOK,
				response:    "{\"id\":\"GetSet186\",\"type\":\"counter\",\"delta\":1}\n",
				contentType: "",
			},
		},
	}
	for _, testUnit := range tests {
		test := testUnit
		t.Run(test.name, func(t *testing.T) {
			echoFramework := echo.New()
			var body io.Reader
			if test.reqJSON != "" {
				body = strings.NewReader(test.reqJSON)
			}

			request := httptest.NewRequest(http.MethodPost, test.request, body)
			responseRecorder := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			c := echoFramework.NewContext(request, responseRecorder)
			err := handler.NewBaseHandler(nil).UpdateHandlerJSON(c)
			assert.NoError(t, err, "expected no errors")

			res := responseRecorder.Result()
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
