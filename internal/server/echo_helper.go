package server

import (
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/sqldb"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware inits and some middlewares to Echo framework.
func SetupMiddleware(echoFramework *echo.Echo, cfg *config.ServerConfig) {
	echoFramework.Use(middleware2.SubnetChecker(*cfg))

	if cfg.CryptoKey != "" {
		echoFramework.Use(middleware2.RsaAesDecoder())
	}
	// Logging middlewares
	// RequestLoggerWithConfig and BodyDump
	loggerConfig := middleware2.GetRequestLoggerConfig()
	echoFramework.Use(middleware.RequestLoggerWithConfig(loggerConfig))
	echoFramework.Use(middleware2.AuthValidator(*cfg))
	echoFramework.Use(middleware.BodyDump(middleware2.GetBodyLoggerHandler()))

	// Set up a compression middleware
	echoFramework.Use(middleware2.GetGzipMiddlewareConfig())
}

// SetupRouter adds some paths to Echo framework.
func SetupRouter(echoFramework *echo.Echo, conn *sqldb.PgxIface) {
	dbHandler := handler.NewBaseHandler(conn)
	echoFramework.POST("/update/:type/:name/:value", handler.UpdateHandler)
	echoFramework.POST("/updates/", dbHandler.UpdatesHandlerJSON)
	echoFramework.POST("/update/", dbHandler.UpdateHandlerJSON)
	echoFramework.GET("/value/:type/:name", dbHandler.ValueHandler)
	echoFramework.POST("/value/", dbHandler.ValueHandlerJSON)
	echoFramework.GET("/", dbHandler.RootHandler)

	echoFramework.GET("/ping", dbHandler.PingHandler)
}
