package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestLogValuesFunc(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	sugar := *logger.Sugar()
	zapSugar = sugar
	e := echo.New()
	assert.NoError(t, logValuesFunc(e.AcquireContext(), middleware.RequestLoggerValues{}))
}

func TestGetRequestLoggerConfig(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	sugar := *logger.Sugar()

	type args struct {
		sugar zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want middleware.RequestLoggerConfig
	}{
		{
			name: "test RequestLoggerConfig ",
			args: args{
				sugar: sugar,
			},
			want: middleware.RequestLoggerConfig{
				LogURI:           true,
				LogStatus:        true,
				LogLatency:       true,
				LogContentLength: true,
				LogResponseSize:  true,
				LogMethod:        true,
				LogValuesFunc:    logValuesFunc,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRequestLoggerConfig(tt.args.sugar)
			assert.Equal(t, got.LogMethod, tt.want.LogMethod)
			assert.Equal(t, got.LogURI, tt.want.LogURI)
			assert.Equal(t, got.LogStatus, tt.want.LogStatus)
			assert.Equal(t, got.LogResponseSize, tt.want.LogResponseSize)
			assert.Equal(t, got.LogLatency, tt.want.LogLatency)

		})
	}
}
