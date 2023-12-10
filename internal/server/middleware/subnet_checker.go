package middleware

import (
	"fmt"
	"net/http"
	"net/netip"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/labstack/echo/v4"
)

var errNotTrustedSubnet = echo.NewHTTPError(http.StatusForbidden, "no trusted subnet")

// SubnetChecker checks "X-Real-IP" header and its value.
func SubnetChecker(cfg config.ServerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			if cfg.HasTrustedSubnet() {
				realIP := echoCtx.Request().Header.Get("X-Real-IP")
				if ok, err := isTrusted(cfg.TrustedSubnet, realIP); err != nil || !ok {
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

func isTrusted(subnet, realIP string) (bool, error) {
	network, err := netip.ParsePrefix(subnet) // "172.17.0.0/16"
	if err != nil {
		return false, fmt.Errorf("can't parse network by: %w", err)
	}

	ip, err := netip.ParseAddr(realIP) // "172.17.0.2"
	if err != nil {
		return false, fmt.Errorf("can't parse real IP by: %w", err)
	}

	return network.Contains(ip), nil
}
