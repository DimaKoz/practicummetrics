package sender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
	goccyj "github.com/goccy/go-json"
)

// ParcelsSend sends metrics.
func ParcelsSend(cfg *config.AgentConfig, metrics []model.MetricUnit) {
	if len(metrics) < minimumBatchNumber { // sending one by one
		sendingSingle(resty.New(), cfg, metrics)
	} else { // do batch request
		sendingBatch(cfg, metrics)
	}
}

// sendingBatch sends a batch request.
func sendingBatch(cfg *config.AgentConfig, metrics []model.MetricUnit) {
	targetURL := getMetricsUpdateTargetURL(cfg.Address, endpointParcelsSend)
	request := resty.New().R()
	addHeadersToRequest(request)
	metrcsSending := make([]model.Metrics, 0, len(metrics))
	for _, unit := range metrics {
		emptyMetrics := model.NewEmptyMetrics()
		emptyMetrics.UpdateByMetricUnit(unit)
		metrcsSending = append(metrcsSending, *emptyMetrics)
	}
	body, err := goccyj.Marshal(metrcsSending)
	if err != nil {
		logSendingErr(err)

		return
	}
	if cfg.HashKey != "" {
		appendHashOtherMarshaling(request, cfg.HashKey, body)
	}

	request.SetBody(body)
	if _, err = request.Post(targetURL); err != nil {
		logSendingErr(err)
	}
}

// sendingSingle sends a single request.
func sendingSingle(rClient *resty.Client, cfg *config.AgentConfig, metrics []model.MetricUnit) {
	emptyMetrics := model.NewEmptyMetrics()
	targetURL := getMetricsUpdateTargetURL(cfg.Address, endpointParcelSend)
	for _, unit := range metrics {
		request := rClient.R()
		emptyMetrics.UpdateByMetricUnit(unit)

		body, err := goccyj.Marshal(emptyMetrics)
		if err != nil {
			logSendingErr(err)

			return
		}

		request.SetBody(body)
		if cfg.HashKey != "" {
			appendHashOtherMarshaling(request, cfg.HashKey, body)
		}
		addHeadersToRequest(request)
		if _, err = request.Post(targetURL); err != nil {
			logSendingErr(err)

			break
		}
	}
}

// appendHashOtherMarshaling appends a hash to common.HashKeyHeaderName header of resty.Request.
func appendHashOtherMarshaling(request *resty.Request, hashKey string, body []byte) {
	key := []byte(hashKey)
	h := hmac.New(sha256.New, key)
	h.Write(body)
	hmacString := hex.EncodeToString(h.Sum(nil))

	request.SetHeader(common.HashKeyHeaderName, hmacString)
}

// logSendingErr prints an error.
func logSendingErr(err error) {
	log.Printf("could not create the request: %s \n", err)
	log.Println("waiting for the next tick")
}

// addHeadersToRequest "Content-Type" and "Accept-Encoding" headers to resty.Request.
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

// getMetricsUpdateTargetURL prepares an URL from address and endpoint.
func getMetricsUpdateTargetURL(address string, endpoint string) string {
	buffLen := len(protocolParcelsSend) + len(endpoint) + len(address)
	strBld := strings.Builder{}
	strBld.Grow(buffLen)
	strBld.WriteString(protocolParcelsSend)
	strBld.WriteString(address)
	strBld.WriteString(endpoint)

	return strBld.String()
}
