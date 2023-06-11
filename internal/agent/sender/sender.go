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
	var targetURL string
	if len(metrics) < minimumBatchNumber { // sending one by one
		emptyMetrics := model.NewEmptyMetrics()
		targetURL = getMetricsUpdateTargetURL(cfg.Address, endpointParcelSend)
		for _, unit := range metrics {
			request := client.R()
			emptyMetrics.UpdateByMetricUnit(unit)
			request.SetBody(emptyMetrics)
			addHeadersToRequest(request)
			if _, err := request.Post(targetURL); err != nil {
				logSendingErr(err)

				break
			}
		}
	} else { // do batch request
		targetURL = getMetricsUpdateTargetURL(cfg.Address, endpointParcelsSend)
		request := client.R()
		addHeadersToRequest(request)
		metrcsSending := make([]model.Metrics, 0, len(metrics))
		for _, unit := range metrics {
			emptyMetrics := model.NewEmptyMetrics()
			emptyMetrics.UpdateByMetricUnit(unit)
			metrcsSending = append(metrcsSending, *emptyMetrics)
		}
		request.SetBody(metrcsSending)
		if _, err := request.Post(targetURL); err != nil {
			logSendingErr(err)
		}
	}
}

func logSendingErr(err error) {
	log.Printf("could not create the request: %s \n", err)
	log.Println("waiting for the next tick")
}

func addHeadersToRequest(request *resty.Request) {
	request.SetHeader("Content-Type", "application/json")
	request.SetHeader("Accept-Encoding", "gzip")
}

const (
	protocolParcelsSend = "http://"
	endpointParcelSend  = "/update/"
	endpointParcelsSend = "/updates/"
	minimumBatchNumber  = 2
)

func getMetricsUpdateTargetURL(address string, endpoint string) string {
	buffLen := len(protocolParcelsSend) + len(endpoint) + len(address)
	strBld := strings.Builder{}
	strBld.Grow(buffLen)
	strBld.WriteString(protocolParcelsSend)
	strBld.WriteString(address)
	strBld.WriteString(endpoint)

	return strBld.String()
}
