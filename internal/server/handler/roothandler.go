package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

// RootHandler handles `/`
func RootHandler(c echo.Context) error {
	metrics := repository.GetAllMetrics()
	str := getHTMLContent(metrics)
	c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
	return c.String(http.StatusOK, str)
}

func getHTMLContent(metrics []model.MetricUnit) string {
	var body = ""
	for i, m := range metrics {
		if i != 0 {
			body += "<br></br>"
		}
		body += m.Name + "," + m.Value
	}

	return fmt.Sprintf("<h1>%s</h1><div>%s</div>", "Metrics:", body)
}
