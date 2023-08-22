package server

import (
	"fmt"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	middleware2 "github.com/DimaKoz/practicummetrics/internal/server/middleware"
	goccyj "github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware inits and some middlewares to Echo framework.
func SetupMiddleware(echoFramework *echo.Echo, cfg *config.ServerConfig) {
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
func SetupRouter(echoFramework *echo.Echo, conn *pgx.Conn) {
	dbHandler := handler.NewBaseHandler(conn)
	echoFramework.POST("/update/:type/:name/:value", handler.UpdateHandler)
	echoFramework.POST("/updates/", dbHandler.UpdatesHandlerJSON)
	echoFramework.POST("/update/", dbHandler.UpdateHandlerJSON)
	echoFramework.GET("/value/:type/:name", dbHandler.ValueHandler)
	echoFramework.POST("/value/", dbHandler.ValueHandlerJSON)
	echoFramework.GET("/", dbHandler.RootHandler)

	echoFramework.GET("/ping", dbHandler.PingHandler)

	// pprof.Register(echoFramework)
}

// FastJSONSerializer implements JSON encoding using encoding/json.
type FastJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d FastJSONSerializer) Serialize(c echo.Context, data interface{}, indent string) error {
	enc := goccyj.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(data) //nolint:wrapcheck
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d FastJSONSerializer) Deserialize(c echo.Context, data interface{}) error {
	err := goccyj.NewDecoder(c.Request().Body).Decode(data)
	if ute, ok := err.(*goccyj.UnmarshalTypeError); ok { //nolint:errorlint
		mess := fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v",
			ute.Type, ute.Value, ute.Field, ute.Offset)

		return echo.NewHTTPError(http.StatusBadRequest, mess).SetInternal(err)
	} else if syne, ok := err.(*goccyj.SyntaxError); ok { //nolint:errorlint
		mess := fmt.Sprintf("Syntax error: offset=%v, error=%v", syne.Offset, syne.Error())

		return echo.NewHTTPError(http.StatusBadRequest, mess).SetInternal(err)
	}

	return err //nolint:wrapcheck
}
