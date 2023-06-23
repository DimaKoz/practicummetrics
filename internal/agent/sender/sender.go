package sender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
)

// ParcelsSend sends metrics.
func ParcelsSend(cfg *config.AgentConfig, metrics []model.MetricUnit) {
	if len(metrics) < minimumBatchNumber { // sending one by one
		sendingSingle(resty.New(), cfg, metrics)
	} else { // do batch request
		sendingBatch(cfg, metrics)
	}
}

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
	if cfg.HashKey != "" {
		if err := appendHash(request, cfg.HashKey, metrcsSending); err != nil {
			logSendingErr(err)

			return
		}
	}

	request.SetBody(metrcsSending)
	if _, err := request.Post(targetURL); err != nil {
		logSendingErr(err)
	}
}

func sendingSingle(rClient *resty.Client, cfg *config.AgentConfig, metrics []model.MetricUnit) {
	emptyMetrics := model.NewEmptyMetrics()
	targetURL := getMetricsUpdateTargetURL(cfg.Address, endpointParcelSend)
	for _, unit := range metrics {
		request := rClient.R()
		emptyMetrics.UpdateByMetricUnit(unit)
		request.SetBody(emptyMetrics)
		if cfg.HashKey != "" {
			if err := appendHash(request, cfg.HashKey, emptyMetrics); err != nil {
				logSendingErr(err)

				return
			}
		}
		addHeadersToRequest(request)
		if _, err := request.Post(targetURL); err != nil {
			logSendingErr(err)

			break
		}
	}
}

func appendHash(request *resty.Request, hashKey string, v interface{}) error {
	bJSON, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("agent: can't get json: %w", err)
	}

	key := []byte(hashKey)
	h := hmac.New(sha256.New, key)
	h.Write(bJSON)
	hmacString := hex.EncodeToString(h.Sum(nil))

	request.SetHeader(common.HashKeyHeaderName, hmacString)

	return nil
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
