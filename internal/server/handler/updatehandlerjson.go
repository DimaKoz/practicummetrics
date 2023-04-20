package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

// UpdateHandlerJSON handles `/update` with json
func UpdateHandlerJSON(c echo.Context) error {
	log.Println("UpdateHandlerJSON")
	m := &model.Metrics{}
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("cannot parse from json: %s", err))
	}

	var metricValue string
	if m.MType == model.MetricTypeGauge {
		metricValue = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	} else {
		metricValue = strconv.FormatInt(*m.Delta, 10)
	}
	muIncome, err := model.NewMetricUnit(m.MType, m.ID, metricValue)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == model.ErrorUnknownType {
			statusCode = http.StatusNotImplemented
		}
		return c.String(statusCode, fmt.Sprintf("cannot create metric: %s", err))
	}
	repository.AddMetric(muIncome)
	if mu, err := repository.GetMetricByName(muIncome.Name); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("cannot find updated metric: %s", err))
	} else {
		m.Convert(mu)
	}

	return c.JSON(http.StatusOK, m)
}
