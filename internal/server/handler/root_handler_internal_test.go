package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(test.method, test.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			_ = RootHandler(c)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode, "StatusCode got: %v, want: %v", res.StatusCode, test.want.code)

			_ = res.Body.Close()
		})
	}
}

func TestGetHtmlContent(t *testing.T) {
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
				{Type: model.MetricTypeCounter, Name: "test", Value: "42", ValueInt: 42, ValueFloat: 0},
				{Type: model.MetricTypeCounter, Name: "test2", Value: "22", ValueInt: 22, ValueFloat: 0},
			},
			want: "<h1>Metrics:</h1><div>test,42<br></br>test2,22</div>",
		},
	}

	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			assert.Equalf(t, test.want, getHTMLContent(test.metrics), "getHtmlContent(%v)", test.metrics)
		})
	}
}
