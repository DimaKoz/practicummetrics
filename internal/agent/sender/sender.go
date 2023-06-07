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
	emptyMetrics := model.NewEmptyMetrics()

	for _, unit := range metrics {
		request := client.R()
		request.SetHeader("Content-Type", "application/json")
		request.SetHeader("Accept-Encoding", "gzip")
		emptyMetrics.UpdateByMetricUnit(unit)
		request.SetBody(emptyMetrics)

		if _, err := request.Post(targetURL); err != nil {
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
	strBld := strings.Builder{}
	strBld.Grow(buffLen)
	strBld.WriteString(protocolParcelsSend)
	strBld.WriteString(address)
	strBld.WriteString(endpointParcelsSend)

	return strBld.String()
}
