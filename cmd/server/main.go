package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/common/sqldb"
	"github.com/DimaKoz/practicummetrics/internal/server"
	"github.com/DimaKoz/practicummetrics/internal/server/grpcsrv"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/DimaKoz/practicummetrics/internal/server/serializer"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

// DB connection
// urlExample := "postgres://videos:userpassword@localhost:5432/testdb"
// urlExample := "postgres://localhost:5432/testdb?sslmode=disable"
// _ = os.Setenv("DATABASE_DSN", urlExample)

func main() {
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
	zap.S().Infoln(config.PrepBuildValues(BuildVersion, BuildDate, BuildCommit))

	printCfgInfo(cfg)

	repository.LoadPrivateKey(*cfg)
	ctx := context.Background()
	var pgxConn sqldb.PgxIface
	var conn *sqldb.PgxIface
	if pgxConn, err = sqldb.ConnectDB(cfg); err == nil {
		defer pgxConn.Close(ctx)
		conn = &pgxConn
	} else {
		conn = nil
		zap.S().Warnf("failed to get a db connection by %s", err.Error())
	}

	if _, err = os.Stat(cfg.FileStoragePath); os.IsNotExist(err) {
		cfg.Restore = false
		zap.S().Info("%v file does not exist\n", cfg.FileStoragePath)
	}

	repository.SetupFilePathStorage(cfg.FileStoragePath)
	loadIfNeed(cfg)
	sGrpc, err := grpcsrv.New(*cfg)
	if err != nil {
		zap.S().Fatal(err)
	}

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
	if err = sGrpc.Run(ctx); err != nil {
		zap.S().Fatal(fmt.Errorf("starting grpc server failed by: %w", err))
	}
	startServer(cfg, conn)
}

func loadIfNeed(cfg *config.ServerConfig) {
	if cfg.IsUseDatabase() {
		return
	}
	needLoad := cfg.Restore && cfg.FileStoragePath != ""
	if needLoad {
		if err := repository.LoadVariant(); err != nil {
			zap.S().Errorf("couldn't restore metrics by %s", err)
		}
	}
}

func printCfgInfo(cfg *config.ServerConfig) {
	zap.S().Info(
		"cfg: \n", cfg.StringVariantCopy(),
	)
	zap.S().Infow("Starting server")
}

func startServer(cfg *config.ServerConfig, conn *sqldb.PgxIface) {
	echos := echo.New()
	echos.JSONSerializer = serializer.FastJSONSerializer{}

	server.SetupMiddleware(echos, cfg)
	server.SetupRouter(echos, conn)

	go func(cfg config.ServerConfig, echoFramework *echo.Echo) {
		zap.S().Info("start server")
		if err := echoFramework.Start(cfg.Address); err != nil && errors.Is(err, http.ErrServerClosed) {
			zap.S().Info("shutting down the server")
		}
	}(*cfg, echos)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	zap.S().Info("awaiting a signal or press Ctrl+C to finish this server")
	<-quit
	zap.S().Info("quit...")
	timeoutDelay := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutDelay)*time.Second)
	defer cancel()
	if err := echos.Shutdown(ctx); err != nil {
		zap.S().Fatal(err)
	}
}
