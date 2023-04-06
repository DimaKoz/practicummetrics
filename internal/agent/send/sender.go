package send

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"net/http"
	"os"
)

// ParcelsSend sends metrics
func ParcelsSend(metrics []model.MetricUnit) {
	for _, unit := range metrics {
		url := "http://localhost:8080/update/" + prepPathByMetric(unit)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}
		res.Body.Close()
	}

}

func prepPathByMetric(mu model.MetricUnit) string {
	return mu.Type + "/" + mu.Name + "/" + mu.Value
}
