package handler

import (
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ValueHandler handles `/value/`
func ValueHandler(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.String(http.StatusNotFound, "couldn't find a name of a metric")
	}

	mu := repository.GetMetricByName(name)
	if mu == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.String(http.StatusOK, mu.Value)
}
