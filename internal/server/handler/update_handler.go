package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

// UpdateHandler handles `/update/`.
func UpdateHandler(ctx echo.Context) error {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")
	metricValue := ctx.Param("value")
	metricUnit, err := model.NewMetricUnit(metricType, metricName, metricValue)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, model.ErrUnknownType) {
			statusCode = http.StatusNotImplemented
		}

		return ctx.String(statusCode, fmt.Sprintf("cannot create metric: %s", err))
	}

	repository.AddMetric(metricUnit)

	return ctx.NoContent(http.StatusOK)
}
