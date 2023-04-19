package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ValueHandler handles `/value/`
func ValueHandler(c echo.Context) error {
	name := c.Param("name")

	mu, err := repository.GetMetricByName(name)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf(" 'value' handler: %s", err.Error()))
	}

	return c.String(http.StatusOK, mu.Value)
}
