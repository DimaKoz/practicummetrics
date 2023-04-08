package handler

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
		name   string
		method string
		target string
		want   want
	}{
		{
			name:   "test 404",
			method: http.MethodGet,
			target: "/status",
			want: want{
				code:        http.StatusNotFound,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "test 200",
			method: http.MethodGet,
			target: "/",
			want: want{
				code:        http.StatusOK,
				response:    ``,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(test.method, test.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			RootHandler(c)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code, "StatusCode got: %v, want: %v", res.StatusCode, test.want.code)

			res.Body.Close()

		})
	}
}

func Test_getHtmlContent(t *testing.T) {
	type args struct {
		metrics []model.MetricUnit
	}
	tests := []struct {
		name    string
		metrics []model.MetricUnit
		want    string
	}{
		{
			name:    "no metrics",
			metrics: []model.MetricUnit{},
			want:    "<h1>Metrics:</h1><div></div>",
		},
		{
			name: "2 metrics",
			metrics: []model.MetricUnit{
				{Type: model.MetricTypeCounter, Name: "test", Value: "42", ValueI: 42, ValueF: 0},
				{Type: model.MetricTypeCounter, Name: "test2", Value: "22", ValueI: 22, ValueF: 0},
			},
			want: "<h1>Metrics:</h1><div>test,42<br></br>test2,22</div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getHtmlContent(tt.metrics), "getHtmlContent(%v)", tt.metrics)
		})
	}
}
