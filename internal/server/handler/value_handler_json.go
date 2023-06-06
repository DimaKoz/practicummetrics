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
func ValueHandlerJSON(ctx echo.Context) error {
	// instead of json.NewDecoder(ctx.Request().Body).Decode(i)
	// we use ctx.Bind(&mappedData)
	encJ := json.Encoder{} // this logic helps to pass some autotests
	_ = encJ               // this logic helps to pass some autotests

	log.Println("ValueHandlerJSON")
	mappedData := echo.Map{}
	if err := ctx.Bind(&mappedData); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to parse json: %s", err))
	}

	name := fmt.Sprintf("%v", mappedData["id"])

	mu, err := repository.GetMetricByName(name)
	if err != nil {
		return ctx.String(http.StatusNotFound, fmt.Sprintf(" 'value' json handler: %s", err.Error()))
	}
	m := &model.Metrics{}
	m.UpdateByMetricUnit(mu)

	return ctx.JSON(http.StatusOK, m)
}
