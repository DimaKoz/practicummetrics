package handler

import (
	"errors"
	"fmt"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	okPathPartsValue = 4
	indexTypeValue   = 2
	indexNameValue   = 3
)

// ValueHandler handles `/value/`
func ValueHandler(c echo.Context) error {
	name, err := getNameFromPath(c.Request().URL.Path)
	if err != nil {
		errHTTP := echo.NewHTTPError(err.StatusCode, err.Error())
		c.Error(errHTTP)
		return errHTTP
	}
	mu := repository.GetMetricByName(name)
	if mu == nil {
		errHTTP := echo.NewHTTPError(http.StatusNotFound)
		c.Error(errHTTP)
		return errHTTP
	}

	if err2 := c.String(http.StatusOK, mu.Value); err2 != nil {
		errHTTP := echo.NewHTTPError(http.StatusInternalServerError)
		c.Error(errHTTP)
		fmt.Println("error for ValueHandler: ", err2)

		return errHTTP

	}
	return nil
}

func getNameFromPath(path string) (string, *error2.RequestError) {
	if path == "" {
		return "", &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}
	}
	parts := strings.Split(path, "/")
	if len(parts) != okPathPartsValue {
		return "", &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")}
	}
	return parts[indexNameValue], nil
}
