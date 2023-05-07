package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

var SyncSaveUpdateHandlerJSON = false

// UpdateHandlerJSON handles `/update` with json
func UpdateHandlerJSON(c echo.Context) error {
	m := &model.Metrics{}
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("UpdateHandlerJSON: cannot parse from json: %s", err))
	}

	muIncome, err := model.NewMetricUnit(m.MType, m.ID, m.GetPreparedValue())
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == model.ErrorUnknownType {
			statusCode = http.StatusNotImplemented
		}
		return c.String(statusCode, fmt.Sprintf("UpdateHandlerJSON: cannot create metric: %s", err))
	}
	mu := repository.AddMetric(muIncome)
	m.UpdateByMetricUnit(mu)
	if SyncSaveUpdateHandlerJSON {
		go func() {
			err := repository.Save()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	return c.JSON(http.StatusOK, m)
}