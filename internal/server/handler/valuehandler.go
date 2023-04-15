package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ValueHandler handles `/value/`
func ValueHandler(c echo.Context) error {
	fmt.Println(c.ParamNames(), c.ParamValues())
	name := c.Param("name")
	if name == "" {
		return c.String(http.StatusNotFound, "couldn't find a name of a metric")
	}

	mu := repository.GetMetricByName(name)
	if mu == nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err2 := c.String(http.StatusOK, mu.Value); err2 != nil {
		fmt.Println("error for ValueHandler: ", err2)
		return c.NoContent(http.StatusInternalServerError)
	}

	return nil
}
