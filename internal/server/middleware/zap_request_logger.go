package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var zapSugar zap.SugaredLogger

var logValuesFunc = func(c echo.Context, v middleware.RequestLoggerValues) error {
	zapSugar.Infow("request",
		zap.String("Method", v.Method),
		zap.String("URI", v.URI),
		zap.Duration("latency", v.Latency),
	)
	zapSugar.Infow("response",
		zap.Int("status", v.Status),
		zap.String("length", v.ContentLength),
		zap.Int64("size", v.ResponseSize),
	)

	return nil
}

func GetRequestLoggerConfig(sugar zap.SugaredLogger) middleware.RequestLoggerConfig {
	zapSugar = sugar
	return middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogLatency:       true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogMethod:        true,
		LogValuesFunc:    logValuesFunc,
	}
}

func GetBodyLoggerHandler(sugar zap.SugaredLogger) middleware.BodyDumpHandler {
	return func(c echo.Context, reqBody, resBody []byte) {
		sugar.Infow(
			"body:", "reqBody:", string(reqBody[:]),
		)
	}
}
