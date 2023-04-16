package gather

import (
	"errors"
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var metricsName = []string{
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

var ErrStart = errors.New("error while collecting metrics with: ")

func collectUint64Metrics(rtm *runtime.MemStats) ([]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0, len(metricsName))
	var errR error
	for _, name := range metricsName {
		value := getFieldValue(rtm, name)
		if m, err := model.NewMetricUnit(model.MetricTypeGauge, name, value); err == nil && m != model.EmptyMetric {
			result = append(result, m)
		} else {
			if errR == nil {
				errR = ErrStart
			}
			errR = fmt.Errorf("%w \n can't get '%s' metric by %s ", errR, name, err.Err.Error())
		}
	}
	return result, errR
}

func collectOtherTypeMetrics(rtm *runtime.MemStats) ([]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0)
	var errR error

	// GCCPUFraction
	fraction := strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "GCCPUFraction", fraction); err == nil && m != model.EmptyMetric {
		result = append(result, m)
	} else {
		errR = fmt.Errorf("%w \n can't get '%s' metric by %s ", ErrStart, "GCCPUFraction", err.Err.Error())
	}

	// NumForcedGC
	numForcedGC := strconv.FormatUint(uint64(rtm.NumForcedGC), 10)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "NumForcedGC", numForcedGC); err == nil && m != model.EmptyMetric {
		result = append(result, m)
	} else {
		if errR == nil {
			errR = ErrStart
		}
		errR = fmt.Errorf("%w \n can't get '%s' metric by %s", errR, "NumForcedGC", err.Err.Error())
	}

	// NumGC
	numGC := strconv.FormatUint(uint64(rtm.NumGC), 10)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "NumGC", numGC); err == nil && m != model.EmptyMetric {
		result = append(result, m)
	} else {
		if errR == nil {
			errR = ErrStart
		}
		errR = fmt.Errorf("%w \n can't get '%s' metric by %s", errR, "NumGC", err.Err.Error())
	}

	// RandomValue
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomValue := strconv.Itoa(r1.Intn(100))
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", randomValue); err == nil && m != model.EmptyMetric {
		result = append(result, m)
	} else {
		if errR == nil {
			errR = ErrStart
		}
		errR = fmt.Errorf("%w \n can't get '%s' metric by %s", errR, "NumGC", err.Err.Error())
	}

	// PollCount
	if m, err := model.NewMetricUnit(model.MetricTypeCounter, "PollCount", "1"); err == nil && m != model.EmptyMetric {
		result = append(result, m)
	} else {
		if errR == nil {
			errR = ErrStart
		}
		errR = fmt.Errorf("%w \n can't get '%s' metric by %s", errR, "NumGC", err.Err.Error())
	}

	return result, errR
}

// GetMetrics returns a list of the metrics
func GetMetrics() ([]model.MetricUnit, error) {

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	result, err := collectUint64Metrics(&rtm)

	m, err2 := collectOtherTypeMetrics(&rtm)
	result = append(result, m...)
	if err2 != nil {
		if err == nil {
			err = err2
		} else {
			err = fmt.Errorf("%w \n %s", err, err2.Error())
		}
	}

	return result, err
}

func getFieldValue(e *runtime.MemStats, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)
	return strconv.FormatUint(f.Uint(), 10)
}
