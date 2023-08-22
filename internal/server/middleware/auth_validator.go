package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
)

// AuthValidator checks "HashSHA256" header and its value.
func AuthValidator(cfg config.ServerConfig) echo.MiddlewareFunc {
	badHash := echo.NewHTTPError(http.StatusBadRequest, "bad hash")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			if true { // hash is temporary disabled
				return next(echoCtx)
			}

			if cfg.HashKey == "" {
				return next(echoCtx)
			}

			// Hash key
			headerHash := echoCtx.Request().Header.Get(common.HashKeyHeaderName)
			if headerHash == "" {
				zap.S().Info("ups:", "missed HashSHA256")

				return badHash
			}
			// Request
			reqBody := []byte{}
			if echoCtx.Request().Body != nil { // Read
				reqBody, _ = io.ReadAll(echoCtx.Request().Body)
			}
			echoCtx.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			if isBadHash(cfg.HashKey, headerHash, reqBody) {
				return badHash
			}

			if err := next(echoCtx); err != nil {
				echoCtx.Error(err)
			}

			return nil
		}
	}
}

// isBadHash returns true when 'incomeHash' is wrong.
func isBadHash(cfgKey string, incomeHash string, reqBody []byte) bool {
	key := []byte(cfgKey)
	h := hmac.New(sha256.New, key)
	h.Write(reqBody)
	hmacString := hex.EncodeToString(h.Sum(nil))
	zap.S().Infow(
		"HashSHA256:", " server:", hmacString,
		" agent:", incomeHash,
	)

	return incomeHash != hmacString
}
