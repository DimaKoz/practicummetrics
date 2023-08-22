package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

var syncSaveUpdateHandlerJSON = false

// SetSyncSaveUpdateHandlerJSON enables synchronization to save data to a file.
func SetSyncSaveUpdateHandlerJSON(sync bool) {
	syncSaveUpdateHandlerJSON = sync
}

// UpdateHandlerJSON handles `/update` with json.
func (h *BaseHandler) UpdateHandlerJSON(ctx echo.Context) error {
	metrics := model.NewEmptyMetrics()
	if err := ctx.Bind(&metrics); err != nil {
		return wrapUpdHandlerErr(ctx, http.StatusBadRequest, "UpdateHandlerJSON: failed to parse json: %s", err)
	}
	prepModelValue, err := metrics.GetPreparedValue()
	if err != nil {
		erDesc := fmt.Sprintf("UpdateHandlerJSON: Metrics contains nil: %s", err)

		return wrapUpdHandlerErr(ctx, http.StatusBadRequest, erDesc, err)
	}
	muIncome, err := model.NewMetricUnit(metrics.MType, metrics.ID, prepModelValue)
	if err != nil {
		statusCode := getUpdatesStatusCode(err)

		return wrapUpdHandlerErr(ctx, statusCode, "UpdateHandlerJSON: cannot create metric: %s", err)
	}
	var metricUnit model.MetricUnit
	if h != nil && h.conn != nil {
		if metricUnit, err = repository.AddMetricToDB(h.conn, muIncome); err != nil {
			return wrapUpdHandlerErr(ctx, http.StatusInternalServerError, "UpdateHandlerJSON: cannot create metric: %s", err)
		}
	} else {
		metricUnit = repository.AddMetric(muIncome)
	}

	metrics.UpdateByMetricUnit(metricUnit)

	save()
	if err = ctx.JSON(http.StatusOK, metrics); err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}

func save() {
	if syncSaveUpdateHandlerJSON {
		go func() {
			err := repository.SaveVariant()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func wrapUpdHandlerErr(ctx echo.Context, statusCode int, msg string, errIn error) error {
	err := ctx.String(statusCode, fmt.Sprintf(msg, errIn))
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
