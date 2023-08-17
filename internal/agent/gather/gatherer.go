package gather

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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

func collectUintMetricsVariant(rtm *runtime.MemStats, result *[]model.MetricUnit) error {
	for _, name := range metricsName {
		value := getFieldValue(rtm, name)
		if m, err := model.NewMetricUnit(model.MetricTypeGauge, name, value); err == nil {
			*result = append(*result, m)
		} else {
			return fmt.Errorf(errFormatString, name, err)
		}
	}

	return nil
}

func collectOtherTypeMetricsVariant(rtm *runtime.MemStats, result *[]model.MetricUnit) error {
	// GCCPUFraction
	fraction := strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "GCCPUFraction", fraction); err == nil {
		*result = append(*result, m)
	} else {
		return fmt.Errorf(errFormatString, "GCCPUFraction", err)
	}

	// RandomValue
	var randInterval int64 = 100
	var randomValue string
	if n, err := rand.Int(rand.Reader, big.NewInt(randInterval)); err == nil {
		randomValue = n.String()
	} else {
		return fmt.Errorf(errFormatString, "RandomValue", err)
	}

	if m, err := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", randomValue); err == nil {
		*result = append(*result, m)
	} else {
		return fmt.Errorf(errFormatString, "RandomValue", err)
	}

	// PollCount
	if m, err := model.NewMetricUnit(model.MetricTypeCounter, "PollCount", "1"); err == nil {
		*result = append(*result, m)
	} else {
		return fmt.Errorf(errFormatString, "PollCount", err)
	}

	return nil
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
	var randInterval int64 = 100
	var randomValue string
	if n, err := rand.Int(rand.Reader, big.NewInt(randInterval)); err == nil {
		randomValue = n.String()
	} else {
		return nil, fmt.Errorf(errFormatString, "RandomValue", err)
	}

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

// GetMemoryMetrics returns a list of the metrics.
func GetMemoryMetrics(resultChan chan *[]model.MetricUnit, errChan chan error) {
	const metricsCount = 3

	result := make([]model.MetricUnit, 0, metricsCount)
	var name string

	// CPUutilization1

	name = "CPUutilization1"

	utilization, err := cpu.Percent(time.Duration(0), false)
	if err != nil {
		errChan <- fmt.Errorf(errFormatString, name, err)

		return
	}

	utilizationStValue := strconv.FormatFloat(utilization[0], 'f', 2, 64)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, name, utilizationStValue); err == nil {
		result = append(result, m)
	} else {
		errChan <- fmt.Errorf(errFormatString, name, err)

		return
	}

	// TotalMemory
	name = "TotalMemory"
	virMem, err := mem.VirtualMemory()
	if err != nil {
		errChan <- fmt.Errorf(errFormatString, name, err)

		return
	}

	if m, err := model.NewMetricUnit(model.MetricTypeGauge,
		name,
		strconv.FormatUint(virMem.Total, 10)); err == nil {
		result = append(result, m)
	} else {
		errChan <- fmt.Errorf(errFormatString, name, err)

		return
	}

	// FreeMemory
	name = "FreeMemory"
	if m, err := model.NewMetricUnit(model.MetricTypeGauge,
		name,
		strconv.FormatUint(virMem.Total, 10)); err == nil {
		result = append(result, m)
	} else {
		errChan <- fmt.Errorf(errFormatString, name, err)

		return
	}

	resultChan <- &result
}

// GetMetricsVariant sends a list of the metrics to resultChan.
func GetMetricsVariant(resultChan chan *[]model.MetricUnit, errChan chan error) {
	const metricsCount = 3
	var (
		rtm    runtime.MemStats
		result = make([]model.MetricUnit, 0, metricsCount+len(metricsName))
		err    error
	)

	runtime.ReadMemStats(&rtm)

	if err = collectUintMetricsVariant(&rtm, &result); err != nil {
		errChan <- fmt.Errorf("cannot collectUintMetrics: %w", err)

		return
	}

	if err = collectOtherTypeMetricsVariant(&rtm, &result); err != nil {
		errChan <- fmt.Errorf("cannot collectOtherTypeMetrics: %w", err)

		return
	}

	resultChan <- &result
}

// GetMetrics returns a list of the metrics.
func GetMetrics(resultChan chan *[]model.MetricUnit, errChan chan error) {
	var (
		rtm                 runtime.MemStats
		result, metricUnits *[]model.MetricUnit
		err                 error
	)

	runtime.ReadMemStats(&rtm)

	if result, err = collectUintMetrics(&rtm); err != nil {
		errChan <- fmt.Errorf("cannot collectUintMetrics: %w", err)

		return
	}

	if metricUnits, err = collectOtherTypeMetrics(&rtm); err != nil {
		errChan <- fmt.Errorf("cannot collectOtherTypeMetrics: %w", err)

		return
	}

	*result = append(*result, *metricUnits...)
	resultChan <- result
}

func getFieldValueUint64(e *runtime.MemStats, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)

	return strconv.FormatUint(f.Uint(), 10)
}

//nolint:cyclop
func getFieldValue(mStat *runtime.MemStats, field string) string { //nolint:funlen
	switch field {
	case "NumForcedGC":
		return strconv.FormatUint(uint64(mStat.NumForcedGC), 10)
	case "NumGC":
		return strconv.FormatUint(uint64(mStat.NumGC), 10)
	case "Alloc":
		return strconv.FormatUint(mStat.Alloc, 10)
	case "BuckHashSys":
		return strconv.FormatUint(mStat.BuckHashSys, 10)
	case "Frees":
		return strconv.FormatUint(mStat.Frees, 10)
	case "GCSys":
		return strconv.FormatUint(mStat.GCSys, 10)
	case "HeapAlloc":
		return strconv.FormatUint(mStat.HeapAlloc, 10)
	case "HeapIdle":
		return strconv.FormatUint(mStat.HeapIdle, 10)
	case "HeapInuse":
		return strconv.FormatUint(mStat.HeapInuse, 10)
	case "HeapObjects":
		return strconv.FormatUint(mStat.HeapObjects, 10)
	case "HeapReleased":
		return strconv.FormatUint(mStat.HeapReleased, 10)
	case "HeapSys":
		return strconv.FormatUint(mStat.HeapSys, 10)
	case "LastGC":
		return strconv.FormatUint(mStat.LastGC, 10)
	case "Lookups":
		return strconv.FormatUint(mStat.Lookups, 10)
	case "MCacheInuse":
		return strconv.FormatUint(mStat.MCacheInuse, 10)
	case "MCacheSys":
		return strconv.FormatUint(mStat.MCacheSys, 10)
	case "MSpanInuse":
		return strconv.FormatUint(mStat.MSpanInuse, 10)
	case "MSpanSys":
		return strconv.FormatUint(mStat.MSpanSys, 10)
	case "Mallocs":
		return strconv.FormatUint(mStat.Mallocs, 10)
	case "NextGC":
		return strconv.FormatUint(mStat.NextGC, 10)
	case "OtherSys":
		return strconv.FormatUint(mStat.OtherSys, 10)
	case "PauseTotalNs":
		return strconv.FormatUint(mStat.PauseTotalNs, 10)
	case "StackInuse":
		return strconv.FormatUint(mStat.StackInuse, 10)
	case "StackSys":
		return strconv.FormatUint(mStat.StackSys, 10)
	case "Sys":
		return strconv.FormatUint(mStat.Sys, 10)
	case "TotalAlloc":
		return strconv.FormatUint(mStat.TotalAlloc, 10)
	}

	return ""
}
