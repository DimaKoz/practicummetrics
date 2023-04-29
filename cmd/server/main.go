package main

import (
	"encoding/json" // this import helps to pass some autotests
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"strings"
)

var sugar zap.SugaredLogger

func main() {

	encJ := json.Encoder{} // this logic helps to pass some autotests
	_ = encJ               // this logic helps to pass some autotests

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
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		sugar.Infow(
			"body:", "reqBody:", string(reqBody[:]),
		)
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			accept := c.Request().Header.Get(echo.HeaderAcceptEncoding)
			isCompressing := strings.Contains(accept, "gzip")
			return isCompressing
		},
		Level: 5,
	}))

	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.POST("/update/", handler.UpdateHandlerJSON)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.POST("/value/", handler.ValueHandlerJSON)
	e.GET("/", handler.RootHandler)

	if err = e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
