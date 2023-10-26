package main

import (
	"context"
	"os"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/server"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/DimaKoz/practicummetrics/internal/server/sqldb"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	// DB connection
	// urlExample := "postgres://videos:userpassword@localhost:5432/testdb"
	// urlExample := "postgres://localhost:5432/testdb?sslmode=disable"
	// _ = os.Setenv("DATABASE_DSN", urlExample)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer func(loggerZap *zap.Logger) {
		_ = loggerZap.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)

	cfg := config.NewServerConfig()

	if err = config.LoadServerConfig(cfg, config.ProcessEnvServer); err != nil {
		zap.S().Fatalf("couldn't create a config %s", err)
	}

	printCfgInfo(cfg)

	var conn *pgx.Conn
	if conn, err = sqldb.ConnectDB(cfg); err == nil {
		defer conn.Close(context.Background())
	} else {
		zap.S().Warnf("failed to get a db connection by %s", err.Error())
	}

	if _, err = os.Stat(cfg.FileStoragePath); os.IsNotExist(err) {
		cfg.Restore = false
		zap.S().Info("%v file does not exist\n", cfg.FileStoragePath)
	}

	repository.SetupFilePathStorage(cfg.FileStoragePath)
	loadIfNeed(cfg)

	if cfg.FileStoragePath != "" {
		if cfg.StoreInterval != 0 {
			ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

			defer ticker.Stop()

			go func() {
				tickerChannel := ticker.C
				for range tickerChannel {
					if err = repository.SaveVariant(); err != nil {
						zap.S().Fatalf("server: cannot save metrics: %s", err)
					}
				}
			}()
		} else {
			handler.SetSyncSaveUpdateHandlerJSON(true)
		}
	}

	startServer(cfg, conn)
}

func loadIfNeed(cfg *config.ServerConfig) {
	if cfg.IsUseDatabase() {
		return
	}
	needLoad := cfg.Restore && cfg.FileStoragePath != ""
	if needLoad {
		if err := repository.Load(); err != nil {
			zap.S().Fatalf("couldn't restore metrics by %s", err)
		}
	}
}

func printCfgInfo(cfg *config.ServerConfig) {
	zap.S().Info(
		"cfg: \n", cfg.StringVariantCopy(),
	)
	zap.S().Infow("Starting server")
}

func startServer(cfg *config.ServerConfig, conn *pgx.Conn) {
	echos := echo.New()
	echos.JSONSerializer = server.FastJSONSerializer{}

	server.SetupMiddleware(echos, cfg)
	server.SetupRouter(echos, conn)

	if err := echos.Start(cfg.Address); err != nil {
		zap.S().Fatalf("couldn't start the server by %s", err)
	}
}
