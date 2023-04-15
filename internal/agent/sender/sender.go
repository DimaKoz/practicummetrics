package sender

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
	"net/url"
)

// ParcelsSend sends metrics
func ParcelsSend(cfg *config.Config, metrics []model.MetricUnit) {
	for _, unit := range metrics {
		preparedURL := getURL(cfg, unit)
		client := resty.New()
		_, err := client.R().Post(preparedURL.String())
		if err != nil {
			fmt.Printf("client: could not create the request: %s \n", err)
			fmt.Printf("client: waiting for the next tick\n")
			break
		}

	}

}

func getURL(cfg *config.Config, mu model.MetricUnit) url.URL {
	u := url.URL{
		Scheme: "http",
		Host:   cfg.Address,
		Path:   mu.GetPath(),
	}
	return u
}
