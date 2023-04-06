package send

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
	"os"
)

// ParcelsSend sends metrics
func ParcelsSend(metrics []model.MetricUnit) {
	for _, unit := range metrics {
		url := "http://localhost:8080/update/" + prepPathByMetric(unit)
		client := resty.New()
		_, err := client.R().Post(url)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

	}

}

func prepPathByMetric(mu model.MetricUnit) string {
	return mu.Type + "/" + mu.Name + "/" + mu.Value
}
