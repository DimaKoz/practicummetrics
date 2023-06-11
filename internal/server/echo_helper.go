package server

import (
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// SetupMiddleware inits and some middlewares to Echo framework.
func SetupMiddleware(echoFramework *echo.Echo, logger zap.SugaredLogger) {
	// Logging middlewares
	// RequestLoggerWithConfig and BodyDump
	loggerConfig := middleware2.GetRequestLoggerConfig(logger)
	echoFramework.Use(middleware.RequestLoggerWithConfig(loggerConfig))
	echoFramework.Use(middleware.BodyDump(middleware2.GetBodyLoggerHandler(logger)))

	// Set up a compression middleware
	echoFramework.Use(middleware2.GetGzipMiddlewareConfig())
}

// SetupRouter adds some paths to Echo framework.
func SetupRouter(echoFramework *echo.Echo, conn *pgx.Conn) {
	dbHandler := handler.NewBaseHandler(conn)
	echoFramework.POST("/update/:type/:name/:value", handler.UpdateHandler)
	echoFramework.POST("/updates/", dbHandler.UpdatesHandlerJSON)
	echoFramework.POST("/update/", dbHandler.UpdateHandlerJSON)
	echoFramework.GET("/value/:type/:name", dbHandler.ValueHandler)
	echoFramework.POST("/value/", dbHandler.ValueHandlerJSON)
	echoFramework.GET("/", dbHandler.RootHandler)

	echoFramework.GET("/ping", dbHandler.PingHandler)
}
