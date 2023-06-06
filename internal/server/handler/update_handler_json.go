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
	m := &model.Metrics{}
	if err := ctx.Bind(&m); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("UpdateHandlerJSON: failed to parse json: %s", err))
	}
	prepModelValue, err := m.GetPreparedValue()
	if err != nil {
		erDesc := fmt.Sprintf("UpdateHandlerJSON: Metrics contains nil: %s", err)
		if err = ctx.String(http.StatusBadRequest, erDesc); err != nil {
			err = fmt.Errorf("failed to update metrics: %w", err)
		}

		return err
	}
	muIncome, err := model.NewMetricUnit(m.MType, m.ID, prepModelValue)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, model.ErrUnknownType) {
			statusCode = http.StatusNotImplemented
		}

		return ctx.String(statusCode, fmt.Sprintf("UpdateHandlerJSON: cannot create metric: %s", err))
	}
	mu := repository.AddMetric(muIncome)
	m.UpdateByMetricUnit(mu)

	if syncSaveUpdateHandlerJSON {
		go func() {
			err := repository.Save()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	return ctx.JSON(http.StatusOK, m)
}
