package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
)

// UpdatesHandlerJSON handles `/updates/` with json.
func (h *BaseHandler) UpdatesHandlerJSON(ctx echo.Context) error {
	metricsSlice := make([]model.Metrics, 0)

	if err := ctx.Bind(&metricsSlice); err != nil {
		return wrapUpdsHandlerErr(ctx, http.StatusBadRequest, "UpdatesHandlerJSON: failed to parse json: %s", err)
	}
	metricUnits := make([]model.MetricUnit, 0)
	for _, item := range metricsSlice {
		prepModelValue, err := item.GetPreparedValue()
		if err != nil {
			erDesc := fmt.Sprintf("UpdateHandlerJSON: Metrics contains nil: %s", err)

			return wrapUpdsHandlerErr(ctx, http.StatusBadRequest, erDesc, err)
		}
		muIncome, err := model.NewMetricUnit(item.MType, item.ID, prepModelValue)
		if err != nil {
			statusCode := getUpdatesStatusCode(err)

			return wrapUpdsHandlerErr(ctx, statusCode, "UpdatesHandlerJSON: cannot create metric: %s", err)
		}
		metricUnits = append(metricUnits, muIncome)
	}

	if h != nil && h.conn != nil {
		if err := processMetricUnits(ctx, h.conn, metricUnits); err != nil {
			return err
		}
	} else {
		for _, unit := range metricUnits {
			_ = repository.AddMetric(unit)
		}
	}
	saveUpdates()
	if err := ctx.NoContent(http.StatusOK); err != nil {
		err = fmt.Errorf("%w", err)

		return err
	}

	return nil
}

func saveUpdates() {
	if syncSaveUpdateHandlerJSON {
		go func() {
			err := repository.SaveVariant()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func getUpdatesStatusCode(err error) int {
	statusCode := http.StatusBadRequest
	if errors.Is(err, model.ErrUnknownType) {
		statusCode = http.StatusNotImplemented
	}

	return statusCode
}

func wrapUpdsHandlerErr(ctx echo.Context, statusCode int, msg string, errIn error) error {
	err := ctx.String(statusCode, fmt.Sprintf(msg, errIn))
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}

func processMetricUnits(ctx echo.Context, conn *pgx.Conn, metricUnits []model.MetricUnit) error {
	transaction, err := conn.Begin(context.TODO())
	if err != nil {
		return wrapUpdsHandlerErr(ctx, http.StatusBadRequest, "UpdatesHandlerJSON: failed to get a transaction: %s", err)
	}

	for _, unit := range metricUnits {
		if _, err := repository.AddMetricTxToDB(&transaction, unit); err != nil {
			_ = transaction.Rollback(context.TODO())

			return wrapUpdsHandlerErr(ctx, http.StatusInternalServerError, "UpdatesHandlerJSON: cannot create metric: %s", err)
		}
	}
	if err = transaction.Commit(context.TODO()); err != nil {
		return wrapUpdsHandlerErr(ctx, http.StatusInternalServerError,
			"UpdatesHandlerJSON: failed to commit a transaction: %s", err)
	}

	return nil
}
