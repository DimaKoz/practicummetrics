package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

var syncSaveUpdateHandlerJSON = false

func SetSyncSaveUpdateHandlerJSON(sync bool) {
	syncSaveUpdateHandlerJSON = sync
}

// UpdateHandlerJSON handles `/update` with json.
func UpdateHandlerJSON(ctx echo.Context) error {
	metrics := model.NewEmptyMetrics()
	if err := ctx.Bind(&metrics); err != nil {
		err = ctx.String(http.StatusBadRequest, fmt.Sprintf("UpdateHandlerJSON: failed to parse json: %s", err))
		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return err
	}
	prepModelValue, err := metrics.GetPreparedValue()
	if err != nil {
		erDesc := fmt.Sprintf("UpdateHandlerJSON: Metrics contains nil: %s", err)

		return wrapUpdHandlerErr(ctx, http.StatusBadRequest, erDesc, err)
	}
	muIncome, err := model.NewMetricUnit(metrics.MType, metrics.ID, prepModelValue)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, model.ErrUnknownType) {
			statusCode = http.StatusNotImplemented
		}

		return wrapUpdHandlerErr(ctx, statusCode, "UpdateHandlerJSON: cannot create metric: %s", err)
	}
	mu := repository.AddMetric(muIncome)
	metrics.UpdateByMetricUnit(mu)

	if syncSaveUpdateHandlerJSON {
		go func() {
			err := repository.Save()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	if err = ctx.JSON(http.StatusOK, metrics); err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}

func wrapUpdHandlerErr(ctx echo.Context, statusCode int, msg string, errIn error) error {
	err := ctx.String(statusCode, fmt.Sprintf(msg, errIn))
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
