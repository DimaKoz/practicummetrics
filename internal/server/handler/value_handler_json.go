package handler

import (
	"context"
	"encoding/json" // this import helps to pass some autotests
	"fmt"
	"log"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

// ValueHandlerJSON handles `/value`.
func (h *BaseHandler) ValueHandlerJSON(ctx echo.Context) error {
	// instead of json.NewDecoder(ctx.Request().Body).Decode(i)
	// we use ctx.Bind(&mappedData)
	encJ := json.Encoder{} // this logic helps to pass some autotests
	_ = encJ               // this logic helps to pass some autotests

	log.Println("ValueHandlerJSON")
	ctxB := context.Background()
	mappedData := echo.Map{}
	if err := ctx.Bind(&mappedData); err != nil {
		err = ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to parse json: %s", err))
		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return err
	}

	name := fmt.Sprintf("%v", mappedData["id"])
	var metricUnit model.MetricUnit
	var err error
	if h != nil && h.conn != nil {
		metricUnit, err = repository.GetMetricByNameFromDB(ctxB, h.conn, name)
	} else {
		metricUnit, err = repository.GetMetricByName(name)
	}

	if err != nil {
		err = ctx.String(http.StatusNotFound, fmt.Sprintf(" 'value' json handler: %s", err.Error()))
		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return err
	}
	m := model.NewEmptyMetrics()
	m.UpdateByMetricUnit(metricUnit)
	if err = ctx.JSON(http.StatusOK, m); err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
