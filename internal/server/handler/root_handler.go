package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
)

// RootHandler handles `/`.
func (h *BaseHandler) RootHandler(ctx echo.Context) error {
	var metrics []model.MetricUnit
	var err error
	if h != nil && h.conn != nil {
		metrics, _ = repository.GetAllMetricsFromDB(context.Background(), h.conn)
	} else {
		metrics = repository.GetAllMetrics()
	}

	str := getHTMLContent(metrics)

	ctx.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")

	if err = ctx.String(http.StatusOK, str); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// getHTMLContent returns a string with HTML tags.
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
