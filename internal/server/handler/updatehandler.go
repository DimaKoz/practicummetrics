package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UpdateHandler handles `/update/`
func UpdateHandler(c echo.Context) error {

	metricType := c.Param("type")
	metricName := c.Param("name")
	metricValue := c.Param("value")
	mu, err := model.NewMetricUnit(metricType, metricName, metricValue)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == model.ErrorUnknownType {
			statusCode = http.StatusNotImplemented
		}
		return c.String(statusCode, fmt.Sprintf("cannot create metric: %s", err))
	}
	repository.AddMetric(mu)
	return c.NoContent(http.StatusOK)
}
