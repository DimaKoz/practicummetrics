package middleware

import (
	"net/http"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/labstack/echo/v4"
)

var errNotTrustedSubnet = echo.NewHTTPError(http.StatusForbidden, "no trusted subnet")

// SubnetChecker checks "X-Real-IP" header and its value.
func SubnetChecker(cfg config.ServerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			if cfg.HasTrustedSubnet() {
				if !strings.Contains(echoCtx.Request().Header.Get("X-Real-IP"), cfg.TrustedSubnet) {
					return errNotTrustedSubnet
				}
			}

			if err := next(echoCtx); err != nil {
				echoCtx.Error(err)

				return err
			}

			return nil
		}
	}
}
