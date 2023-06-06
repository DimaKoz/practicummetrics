package gather

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var metricsName = []string{
	"NumForcedGC", // uint32
	"NumGC",       // uint32
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

const errFormatString = "error while collecting metrics with: \n can't get '%s' metric by %w "

func collectUintMetrics(rtm *runtime.MemStats) (*[]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0, len(metricsName))

	for _, name := range metricsName {
		value := getFieldValueUint64(rtm, name)
		if m, err := model.NewMetricUnit(model.MetricTypeGauge, name, value); err == nil {
			result = append(result, m)
		} else {
			return nil, fmt.Errorf(errFormatString, name, err)
		}
	}

	return &result, nil
}

func collectOtherTypeMetrics(rtm *runtime.MemStats) (*[]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0)

	// GCCPUFraction
	fraction := strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "GCCPUFraction", fraction); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, "GCCPUFraction", err)
	}

	// RandomValue
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	const randInterval = 100

	randomValue := strconv.Itoa(random.Intn(randInterval))
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", randomValue); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, "RandomValue", err)
	}

	// PollCount
	if m, err := model.NewMetricUnit(model.MetricTypeCounter, "PollCount", "1"); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, "PollCount", err)
	}

	return &result, nil
}

// GetMetrics returns a list of the metrics.
func GetMetrics() (*[]model.MetricUnit, error) {
	var (
		rtm                 runtime.MemStats
		result, metricUnits *[]model.MetricUnit
		err                 error
	)

	runtime.ReadMemStats(&rtm)

	if result, err = collectUintMetrics(&rtm); err != nil {
		return nil, fmt.Errorf("cannot collectUintMetrics: %w", err)
	}

	if metricUnits, err = collectOtherTypeMetrics(&rtm); err != nil {
		return metricUnits, fmt.Errorf("cannot collectOtherTypeMetrics: %w", err)
	}

	*result = append(*result, *metricUnits...)

	return result, err
}

func getFieldValueUint64(e *runtime.MemStats, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)

	return strconv.FormatUint(f.Uint(), 10)
}
