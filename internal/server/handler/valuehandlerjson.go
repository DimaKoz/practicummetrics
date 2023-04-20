package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// ValueHandlerJSON handles `/value`
func ValueHandlerJSON(c echo.Context) error {
	log.Println("ValueHandlerJSON")
	mappedData := echo.Map{}
	if err := c.Bind(&mappedData); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("cannot parse from json: %s", err))
	}

	name := fmt.Sprintf("%v", mappedData["id"])

	mu, err := repository.GetMetricByName(name)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf(" 'value' json handler: %s", err.Error()))
	}
	m := &model.Metrics{}
	m.Convert(mu)

	return c.JSON(http.StatusOK, m)
}
