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
		if err = logger.Sync(); err != nil {
			panic(err)
		}
	}(logger)

	sugar := *logger.Sugar()
	cfg := config.NewServerConfig()

	if err = config.LoadServerConfig(cfg, config.ProcessEnvServer); err != nil {
		sugar.Fatalf("couldn't create a config %s", err)
	}

	printCfgInfo(cfg, sugar)

	if _, err = os.Stat(cfg.FileStoragePath); os.IsNotExist(err) {
		cfg.Restore = false
		sugar.Info("%v file does not exist\n", cfg.FileStoragePath)
	}

	repository.SetupFilePathStorage(cfg.FileStoragePath)
	loadIfNeed(cfg, sugar)

	if cfg.FileStoragePath != "" {
		if cfg.StoreInterval != 0 {
			ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

			defer ticker.Stop()

			go func() {
				tickerChannel := ticker.C
				for range tickerChannel {
					if err = repository.Save(); err != nil {
						sugar.Fatalf("server: cannot save metrics: %s", err)
					}
				}
			}()
		} else {
			handler.SetSyncSaveUpdateHandlerJSON(true)
		}
	}

	startServer(cfg, sugar)
}

func loadIfNeed(cfg *config.ServerConfig, sugar zap.SugaredLogger) {
	needLoad := cfg.Restore && cfg.FileStoragePath != ""
	if needLoad {
		if err := repository.Load(); err != nil {
			sugar.Fatalf("couldn't restore metrics by %s", err)
		}
	}
}

func printCfgInfo(cfg *config.ServerConfig, sugar zap.SugaredLogger) {
	sugar.Info(
		"cfg: \n", cfg.String(),
	)
	sugar.Infow("Starting server")
}

func startServer(cfg *config.ServerConfig, sugar zap.SugaredLogger) {
	e := echo.New()
	server.SetupMiddleware(e, sugar)
	server.SetupRouter(e)

	if err := e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
