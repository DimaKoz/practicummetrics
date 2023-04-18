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

func collectUint64Metrics(rtm *runtime.MemStats) (*[]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0, len(metricsName))
	for _, name := range metricsName {
		value := getFieldValueUint64(rtm, name)
		if m, err := model.NewMetricUnit(model.MetricTypeGauge, name, value); err == nil {
			result = append(result, m)
		} else {
			return nil, fmt.Errorf(errFormatString, name, err.Err)
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
		return nil, fmt.Errorf(errFormatString, "GCCPUFraction", err.Err)
	}

	// RandomValue
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomValue := strconv.Itoa(r1.Intn(100))
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", randomValue); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, "RandomValue", err.Err)
	}

	// PollCount
	if m, err := model.NewMetricUnit(model.MetricTypeCounter, "PollCount", "1"); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, "PollCount", err.Err)
	}

	return &result, nil
}

// GetMetrics returns a list of the metrics
func GetMetrics() (*[]model.MetricUnit, error) {

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	var result, m *[]model.MetricUnit

	var err error
	if result, err = collectUint64Metrics(&rtm); err != nil {
		return nil, err
	}

	if m, err = collectOtherTypeMetrics(&rtm); err != nil {
		return m, err
	}

	*result = append(*result, *m...)

	return result, err
}

func getFieldValueUint64(e *runtime.MemStats, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)
	return strconv.FormatUint(f.Uint(), 10)
}
