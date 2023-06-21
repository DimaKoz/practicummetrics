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

	var conn *pgx.Conn
	if conn, err = sqldb.ConnectDB(cfg, sugar); err == nil {
		defer conn.Close(context.Background())
	} else {
		sugar.Warnf("failed to get a db connection by %s", err.Error())
	}

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

	startServer(cfg, conn, sugar)
}

func loadIfNeed(cfg *config.ServerConfig, sugar zap.SugaredLogger) {
	if cfg.IsUseDatabase() {
		return
	}
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

func startServer(cfg *config.ServerConfig, conn *pgx.Conn, sugar zap.SugaredLogger) {
	e := echo.New()
	server.SetupMiddleware(e, sugar)
	server.SetupRouter(e, conn)

	if err := e.Start(cfg.Address); err != nil {
		sugar.Fatalf("couldn't start the server by %s", err)
	}
}
