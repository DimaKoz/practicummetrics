package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var errBadHash = echo.NewHTTPError(http.StatusBadRequest, "bad hash")

// AuthValidator checks "HashSHA256" header and its value.
func AuthValidator(cfg config.ServerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			if err := authValidate(echoCtx, cfg.HashKey); err != nil {
				return err
			}
			if err := next(echoCtx); err != nil {
				echoCtx.Error(err)
			}

			return nil
		}
	}
}

// authValidate returns err when something is wrong.
func authValidate(echoCtx echo.Context, cfgHashKey string) error {
	if cfgHashKey == "" {
		return nil
	}

	// Hash key
	headerHash := echoCtx.Request().Header.Get(common.HashKeyHeaderName)
	if headerHash == "" {
		zap.S().Info("ups:", "missed HashSHA256")

		return errBadHash
	}
	// Request
	reqBody := getRequestBody(echoCtx)
	if isBadHash(cfgHashKey, headerHash, reqBody) {
		return errBadHash
	}

	return nil
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

func getRequestBody(echoCtx echo.Context) []byte {
	reqBody := []byte{}
	if echoCtx.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(echoCtx.Request().Body)
	}
	echoCtx.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	return reqBody
}
