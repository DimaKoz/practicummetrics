package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
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
		err = ctx.String(statusCode, fmt.Sprintf("cannot create metric: %s", err))
		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return err
	}

	repository.AddMetric(metricUnit)

	if err = ctx.NoContent(http.StatusOK); err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
