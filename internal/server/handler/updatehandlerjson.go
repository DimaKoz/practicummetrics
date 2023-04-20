package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// UpdateHandlerJSON handles `/update` with json
func UpdateHandlerJSON(c echo.Context) error {
	log.Println("UpdateHandlerJSON")
	mappedData := echo.Map{}
	if err := c.Bind(&mappedData); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("cannot parse from json: %s", err))
	}

	metricType := fmt.Sprintf("%v", mappedData["type"])
	metricName := fmt.Sprintf("%v", mappedData["id"])
	var metricValue string
	if metricType == "gauge" {
		metricValue = fmt.Sprintf("%v", mappedData["value"])
	} else {
		metricValue = fmt.Sprintf("%v", mappedData["delta"])
	}

	mu, err := model.NewMetricUnit(metricType, metricName, metricValue)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == model.ErrorUnknownType {
			statusCode = http.StatusNotImplemented
		}
		return c.String(statusCode, fmt.Sprintf("cannot create metric: %s", err))
	}
	repository.AddMetric(mu)
	m := &model.Metrics{}
	m.Convert(mu)

	return c.JSON(http.StatusOK, m)
}
