package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

// RootHandler handles `/`.
func RootHandler(ctx echo.Context) error {
	metrics := repository.GetAllMetrics()

	str := getHTMLContent(metrics)

	ctx.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")

	if err := ctx.String(http.StatusOK, str); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func getHTMLContent(metrics []model.MetricUnit) string {
	strBld := strings.Builder{}
	strBld.WriteString("<h1>Metrics:</h1><div>")
	for i, metricUnit := range metrics {
		if i != 0 {
			strBld.WriteString("<br></br>")
		}

		strBld.WriteString(metricUnit.Name)
		strBld.WriteString(",")
		strBld.WriteString(metricUnit.Value)
	}

	strBld.WriteString("</div>")

	return strBld.String()
}
