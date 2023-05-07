package main

import (
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"time"
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
	sugar.Info(
		"cfg: \n", cfg.String(),
	)
	sugar.Infow(
		"Starting server",
	)
	repository.SetupFilePathStorage(cfg.FileStoragePath)
	if cfg.Restore && cfg.FileStoragePath != "" {
		if err = repository.Load(); err != nil {
			sugar.Fatalf("couldn't restore metrics by %s", err)
		}
	}
	cfg.StoreInterval = 5
	if cfg.FileStoragePath != "" {
		if cfg.StoreInterval != 0 {
			handler.SyncSaveUpdateHandlerJSON = false
			ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
			defer ticker.Stop()
			go func() {
				for {
					select {
					case <-ticker.C:
						err = repository.Save()
						if err != nil {
							sugar.Fatalf("agent: cannot collect metrics: %s", err)
						}
					}
				}
			}()
		} else {
			handler.SyncSaveUpdateHandlerJSON = true
		}
	}

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware2.GetRequestLoggerConfig(sugar)))
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		sugar.Infow(
			"body:", "reqBody:", string(reqBody[:]),
		)
	}))
	e.Use(middleware2.GetGzipMiddlewareConfig())

	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.POST("/update/", handler.UpdateHandlerJSON)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.POST("/value/", handler.ValueHandlerJSON)
	e.GET("/", handler.RootHandler)

	if err = e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
