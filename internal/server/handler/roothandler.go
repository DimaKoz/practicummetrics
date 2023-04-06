package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

// RootHandler handles `/`
func RootHandler(c echo.Context) error {
	if c.Request().URL.Path != "/" {
		errHTTP := echo.NewHTTPError(http.StatusNotFound, "wrong url")
		c.Error(errHTTP)
		return errHTTP
	}
	metrics := repository.GetMetricsMemStorage()
	var body = ""
	for i, m := range metrics {
		if i != 0 {
			body += "<br></br>"
		}
		body += m.Name + "," + m.Value
	}

	str := fmt.Sprintf("<h1>%s</h1><div>%s</div>", "Metrics:", body)
	return c.String(http.StatusOK, str)
}
