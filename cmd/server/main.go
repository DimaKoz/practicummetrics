package main

import (
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

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
	e.Use(middleware.RequestLoggerWithConfig(middleware2.GetRequestLoggerConfig(sugar)))

	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.POST("/update", handler.UpdateHandlerJSON)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.POST("/value", handler.ValueHandlerJSON)
	e.GET("/", handler.RootHandler)

	if err = e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
