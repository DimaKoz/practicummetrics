package handler

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// RootHandler handles `/`
func RootHandler(c echo.Context) error {
	metrics := repository.GetAllMetrics()

	str := getHTMLContent(metrics)
	c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
	return c.String(http.StatusOK, str)
}

func getHTMLContent(metrics []model.MetricUnit) string {
	b := strings.Builder{}
	b.WriteString("<h1>Metrics:</h1><div>")
	for i, m := range metrics {
		if i != 0 {
			b.WriteString("<br></br>")
		}
		b.WriteString(m.Name)
		b.WriteString(",")
		b.WriteString(m.Value)
	}
	b.WriteString("</div>")
	return b.String()
}
