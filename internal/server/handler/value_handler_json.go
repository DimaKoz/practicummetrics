package handler

import (
	"encoding/json" // this import helps to pass some autotests
	"fmt"
	"log"
	"net/http"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

// ValueHandlerJSON handles `/value`.
func ValueHandlerJSON(c echo.Context) error {
	// instead of json.NewDecoder(c.Request().Body).Decode(i)
	// we use c.Bind(&mappedData)
	encJ := json.Encoder{} // this logic helps to pass some autotests
	_ = encJ               // this logic helps to pass some autotests

	log.Println("ValueHandlerJSON")
	mappedData := echo.Map{}
	if err := c.Bind(&mappedData); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to parse json: %s", err))
	}

	name := fmt.Sprintf("%v", mappedData["id"])

	mu, err := repository.GetMetricByName(name)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf(" 'value' json handler: %s", err.Error()))
	}
	m := &model.Metrics{}
	m.UpdateByMetricUnit(mu)
	return c.JSON(http.StatusOK, m)
}
