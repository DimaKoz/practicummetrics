package handler

import (
	"errors"
	"fmt"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	okPathParts = 5
	indexType   = 2
	indexName   = 3
	indexValue  = 4
)

// UpdateHandler handles `/update/`
func UpdateHandler(c echo.Context) error {
	fmt.Println("UpdateHandler", c)
	mu, err := processPath(c.Request().URL.Path)
	if err != nil {
		return c.String(err.StatusCode, err.Error())
	}
	repository.AddMetricMemStorage(*mu)
	return c.NoContent(http.StatusOK)
}

func processPath(path string) (*model.MetricUnit, *error2.RequestError) {

	if path == "" {
		return nil, &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}
	}

	parts := strings.Split(path, "/")
	if len(parts) != okPathParts {
		return nil, &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")}
	}

	mu, err := model.NewMetricUnit(parts[indexType], parts[indexName], parts[indexValue])
	if err != nil {
		return nil, err
	}

	return mu, nil

}
