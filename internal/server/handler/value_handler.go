package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

// ValueHandler handles `/value/`.
func (h *BaseHandler) ValueHandler(ctx echo.Context) error {
	name := ctx.Param("name")
	ctxB := context.Background()
	var (
		metricUnit model.MetricUnit
		err        error
	)

	if h != nil && h.conn != nil {
		metricUnit, err = repository.GetMetricByNameFromDB(ctxB, h.conn, name)
	} else {
		metricUnit, err = repository.GetMetricByName(name)
	}

	if err != nil {
		errDesc := fmt.Sprintf(" 'value' handler: %s", err.Error())
		err = fmt.Errorf("failed to get MetricUnit: %w", ctx.String(http.StatusNotFound, errDesc))

		return err
	}

	if err = ctx.String(http.StatusOK, metricUnit.Value); err != nil {
		err = fmt.Errorf("failed to send MetricUnit: %w", err)
	}

	return err
}
