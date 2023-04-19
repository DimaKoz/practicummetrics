package main

import (
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()

	cfg, err := config.LoadServerConfig()
	if err != nil {
		sugar.Fatalf("couldn't create a config %s", err)
	}

	// from cfg:
	sugar.Infow(
		"cfg:",
		"address", cfg.Address,
	)
	sugar.Infow(
		"Starting server",
	)
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.GET("/", handler.RootHandler)

	if err = e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
