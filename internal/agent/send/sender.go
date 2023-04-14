package send

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
)

// ParcelsSend sends metrics
func ParcelsSend(cfg *config.Config, metrics []model.MetricUnit) {
	for _, unit := range metrics {
		url := "http://" + cfg.Address + "/update/" + prepPathByMetric(unit)
		client := resty.New()
		_, err := client.R().Post(url)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			fmt.Printf("client: waiting for the next tick\n")
			break
		}

	}

}

func prepPathByMetric(mu model.MetricUnit) string {
	return mu.Type + "/" + mu.Name + "/" + mu.Value
}
