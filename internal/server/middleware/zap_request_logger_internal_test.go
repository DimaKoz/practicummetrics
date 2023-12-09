package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func setupLogsCapture() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)

	return zap.New(core), logs
}

func TestLogReqBody(t *testing.T) {
	logger, logs := setupLogsCapture()
	prevL := zap.L()
	defer zap.ReplaceGlobals(prevL)
	zap.ReplaceGlobals(logger)

	wantMessage := "1234"

	logReqBodyImpl(logger.Sugar(), []byte(wantMessage))

	assert.False(t, logs.Len() != 1, "No logs")

	entry := logs.All()[0]
	assert.Equal(t, zap.InfoLevel, entry.Level)
	assert.Equal(t, wantMessage, entry.Context[0].String, "Invalid log entry %v", entry)
}

func TestLogReqBodyFunc(t *testing.T) {
	logger, logs := setupLogsCapture()
	prevL := zap.L()
	defer zap.ReplaceGlobals(prevL)
	zap.ReplaceGlobals(logger)

	wantMessage := "12345"

	handler := GetBodyLoggerHandler()
	e := echo.New()

	handler(e.AcquireContext(), []byte(wantMessage), []byte{})

	assert.False(t, logs.Len() != 1, "No logs")

	entry := logs.All()[0]
	assert.Equal(t, zap.InfoLevel, entry.Level)
	assert.Equal(t, wantMessage, entry.Context[0].String, "Invalid log entry %v", entry)
}

func TestLogValuesFunc(t *testing.T) {
	logger := zap.Must(zap.NewDevelopment())

	defer func(loggerZap *zap.Logger) {
		_ = loggerZap.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)
	e := echo.New()
	assert.NoError(t, logValuesFunc(e.AcquireContext(), middleware.RequestLoggerValues{})) //nolint:exhaustruct
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
			want: middleware.RequestLoggerConfig{ //nolint:exhaustruct
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

	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			got := GetRequestLoggerConfig()
			assert.Equal(t, got.LogMethod, test.want.LogMethod)
			assert.Equal(t, got.LogURI, test.want.LogURI)
			assert.Equal(t, got.LogStatus, test.want.LogStatus)
			assert.Equal(t, got.LogResponseSize, test.want.LogResponseSize)
			assert.Equal(t, got.LogLatency, test.want.LogLatency)
		})
	}
}
