package middleware

import (
	"bytes"
	"errors"
	"io"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

var errRsaAesDecode = errors.New("can't decode message")

func RsaAesDecoder() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			encodedKey := echoCtx.Request().Header.Get(common.RsaEncodedKeyHeaderName)
			if encodedKey == "" {
				return errRsaAesDecode
			}
			var reqBody []byte
			if echoCtx.Request().Body == nil { // Read
				return errRsaAesDecode
			}
			reqBody, err := io.ReadAll(echoCtx.Request().Body)
			if err != nil {
				return errRsaAesDecode
			}
			decodedData, err := repository.DecryptBigMessage(reqBody, encodedKey)
			if err != nil {
				return errRsaAesDecode
			}
			echoCtx.Request().Body = io.NopCloser(bytes.NewReader(decodedData))
			echoCtx.Request().ContentLength = int64(len(decodedData))
			if err = next(echoCtx); err != nil {
				echoCtx.Error(err)

				return err
			}

			return nil
		}
	}
}
