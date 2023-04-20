package sender

import (
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
	"log"
	"net/url"
)

// ParcelsSend sends metrics
func ParcelsSend(cfg *config.AgentConfig, metrics []model.MetricUnit) {
	client := resty.New()
	for _, unit := range metrics {
		//preparedURL := getURL(cfg.Address, unit)
		r := client.R()
		r.SetHeader("Content-Type", "application/json")
		m := &model.Metrics{}
		m.Convert(unit)
		r.SetBody(m)
		_, err := r.Post( /*preparedURL.String()*/ "http://" + cfg.Address + "/update/")
		if err != nil {
			log.Printf("client: could not create the request: %s \n", err)
			log.Println("client: waiting for the next tick")
			break
		}

	}

}

func getURL(address string, mu model.MetricUnit) url.URL {
	u := url.URL{
		Scheme: "http",
		Host:   address,
		Path:   mu.GetPath(),
	}
	return u
}
