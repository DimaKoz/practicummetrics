package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var gzipSkipper = func(c echo.Context) bool {
	accept := c.Request().Header.Get(echo.HeaderAcceptEncoding)
	hasNoGzip := !strings.Contains(accept, "gzip")
	if !hasNoGzip {
		c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
	}

	return hasNoGzip
}

func newGzipConfig(skipper middleware.Skipper) middleware.GzipConfig {
	return middleware.GzipConfig{
		Skipper:   skipper,
		Level:     gzip.BestCompression,
		MinLength: 0,
	}
}

func GetGzipMiddlewareConfig() echo.MiddlewareFunc {
	return middleware.GzipWithConfig(newGzipConfig(gzipSkipper))
}
