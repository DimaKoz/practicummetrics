package server

import (
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
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
func SetupRouter(e *echo.Echo) {
	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.POST("/update/", handler.UpdateHandlerJSON)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.POST("/value/", handler.ValueHandlerJSON)
	e.GET("/", handler.RootHandler)
}
