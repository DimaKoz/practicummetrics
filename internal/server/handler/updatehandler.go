package handler

import (
	"fmt"
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

	if len(c.ParamValues()) != okPathParts {
		return c.String(http.StatusNotFound, "wrong number of the parts of the path")
	}
	mu, err := model.NewMetricUnit(c.ParamValues()[indexType], c.ParamValues()[indexName], c.ParamValues()[indexValue])
	if err != nil {
		return c.String(err.StatusCode, fmt.Sprintf("cannot create metric: %s", err))
	}
	repository.AddMetric(mu)
	return c.NoContent(http.StatusOK)
}
