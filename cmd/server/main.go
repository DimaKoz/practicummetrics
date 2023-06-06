package main

import (
	"os"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/server"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

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

	sugar := *logger.Sugar()

	cfg := config.NewServerConfig()

	err = config.LoadServerConfig(cfg, config.ProcessEnvServer)
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
	_, err = os.Stat(cfg.FileStoragePath)
	if os.IsNotExist(err) {
		cfg.Restore = false
		sugar.Info("%v file does not exist\n", cfg.FileStoragePath)
	}
	repository.SetupFilePathStorage(cfg.FileStoragePath)
	if cfg.Restore && cfg.FileStoragePath != "" {
		if err = repository.Load(); err != nil {
			sugar.Fatalf("couldn't restore metrics by %s", err)
		}
	}

	if cfg.FileStoragePath != "" {
		if cfg.StoreInterval != 0 {
			handler.SetSyncSaveUpdateHandlerJSON(false)
			ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

			defer ticker.Stop()

			go func() {
				tickerChannel := ticker.C
				for range tickerChannel {
					err = repository.Save()
					if err != nil {
						sugar.Fatalf("server: cannot save metrics: %s", err)
					}
				}
			}()
		} else {
			handler.SetSyncSaveUpdateHandlerJSON(true)
		}
	}

	e := echo.New()
	server.SetupMiddleware(e, sugar)
	server.SetupRouter(e)

	if err = e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
