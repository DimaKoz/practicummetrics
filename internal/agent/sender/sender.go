package sender

import (
	"log"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
)

// ParcelsSend sends metrics.
func ParcelsSend(cfg *config.AgentConfig, metrics []model.MetricUnit) {
	client := resty.New()
	targetURL := getTargetURL(cfg.Address)
	m := &model.Metrics{}
	for _, unit := range metrics {
		r := client.R()
		r.SetHeader("Content-Type", "application/json")
		r.SetHeader("Accept-Encoding", "gzip")
		m.UpdateByMetricUnit(unit)
		r.SetBody(m)
		_, err := r.Post(targetURL)
		if err != nil {
			log.Printf("could not create the request: %s \n", err)
			log.Println("waiting for the next tick")

			break
		}
	}
}

const (
	protocolParcelsSend = "http://"
	endpointParcelsSend = "/update/"
)

func getTargetURL(address string) string {
	buffLen := len(protocolParcelsSend) + len(endpointParcelsSend) + len(address)
	b := strings.Builder{}
	b.Grow(buffLen)
	b.WriteString(protocolParcelsSend)
	b.WriteString(address)
	b.WriteString(endpointParcelsSend)
	return b.String()
}
