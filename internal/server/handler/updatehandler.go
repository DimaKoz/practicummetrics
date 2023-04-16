package handler

import (
	"errors"
	"fmt"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	okPathParts = 3
	indexType   = 0
	indexName   = 1
	indexValue  = 2
)

// UpdateHandler handles `/update/`
func UpdateHandler(c echo.Context) error {
	fmt.Println("UpdateHandler", c)
	mu, err := processPath(c, c.Request().URL.Path)
	if err != nil {
		return c.String(err.StatusCode, err.Error())
	}
	repository.AddMetric(mu)
	return c.NoContent(http.StatusOK)
}

func processPath(c echo.Context, path string) (model.MetricUnit, *error2.RequestError) {

	if path == "" {
		return model.MetricUnit{}, &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}
	}
	
	if len(c.ParamValues()) != okPathParts {
		return model.MetricUnit{}, &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")}
	}

	mu, err := model.NewMetricUnit(c.ParamValues()[indexType], c.ParamValues()[indexName], c.ParamValues()[indexValue])
	if err != nil {
		return model.MetricUnit{}, err
	}

	return mu, nil

}
