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
	echoFramework.POST("/update/:type/:name/:value", handler.UpdateHandler)
	echoFramework.POST("/update/", handler.UpdateHandlerJSON)
	echoFramework.GET("/value/:type/:name", handler.ValueHandler)
	echoFramework.POST("/value/", handler.ValueHandlerJSON)
	echoFramework.GET("/", handler.RootHandler)

	dbHandler := handler.NewBaseHandler(conn)
	echoFramework.GET("/ping", dbHandler.PingHandler)
}
