package handler

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"net/http"
)

// RootHandler handles `/`
func RootHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metrics := repository.GetMetricsMemStorage()
	var body = ""
	for i, m := range metrics {
		if i != 0 {
			body += "<br></br>"
		}
		body += m.Name + "," + m.Value
	}

	fmt.Fprintf(res, "<h1>%s</h1><div>%s</div>", "Metrics:", body)

}
