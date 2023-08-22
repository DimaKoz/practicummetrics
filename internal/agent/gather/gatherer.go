package gather

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
)

const (
	MetricNameNumForcedGC   = "NumForcedGC"
	MetricNameNumGC         = "NumGC"
	MetricNameAlloc         = "Alloc"
	MetricNameBuckHashSys   = "BuckHashSys"
	MetricNameFrees         = "Frees"
	MetricNameGCSys         = "GCSys"
	MetricNameHeapAlloc     = "HeapAlloc"
	MetricNameHeapIdle      = "HeapIdle"
	MetricNameHeapInuse     = "HeapInuse"
	MetricNameHeapObjects   = "HeapObjects"
	MetricNameHeapReleased  = "HeapReleased"
	MetricNameHeapSys       = "HeapSys"
	MetricNameLastGC        = "LastGC"
	MetricNameLookups       = "Lookups"
	MetricNameMCacheInuse   = "MCacheInuse"
	MetricNameMCacheSys     = "MCacheSys"
	MetricNameMSpanInuse    = "MSpanInuse"
	MetricNameMSpanSys      = "MSpanSys"
	MetricNameMallocs       = "Mallocs"
	MetricNameNextGC        = "NextGC"
	MetricNameOtherSys      = "OtherSys"
	MetricNamePauseTotalNs  = "PauseTotalNs"
	MetricNameStackInuse    = "StackInuse"
	MetricNameStackSys      = "StackSys"
	MetricNameSys           = "Sys"
	MetricNameTotalAlloc    = "TotalAlloc"
	MetricNameGCCPUFraction = "GCCPUFraction"
	MetricNamePollCount     = "PollCount"
	MetricNameRandomValue   = "RandomValue"
	MetricNameCPUutiliz1    = "CPUutilization1"
	MetricNameTotalMemory   = "TotalMemory"
	MetricNameFreeMemory    = "FreeMemory"
)

var metricsName = []string{
	MetricNameNumForcedGC, // uint32
	MetricNameNumGC,       // uint32
	MetricNameAlloc,
	MetricNameBuckHashSys,
	MetricNameFrees,
	MetricNameGCSys,
	MetricNameHeapAlloc,
	MetricNameHeapIdle,
	MetricNameHeapInuse,
	MetricNameHeapObjects,
	MetricNameHeapReleased,
	MetricNameHeapSys,
	MetricNameLastGC,
	MetricNameLookups,
	MetricNameMCacheInuse,
	MetricNameMCacheSys,
	MetricNameMSpanInuse,
	MetricNameMSpanSys,
	MetricNameMallocs,
	MetricNameNextGC,
	MetricNameOtherSys,
	MetricNamePauseTotalNs,
	MetricNameStackInuse,
	MetricNameStackSys,
	MetricNameSys,
	MetricNameTotalAlloc,
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
		value := getFieldValueVariant(rtm, name)
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
	fraction := fmt.Sprintf("%v", rtm.GCCPUFraction)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, MetricNameGCCPUFraction, fraction); err == nil {
		*result = append(*result, m)
	} else {
		return fmt.Errorf(errFormatString, MetricNameGCCPUFraction, err)
	}

	// RandomValue
	var randInterval int64 = 100
	var randomValue string
	if n, err := rand.Int(rand.Reader, big.NewInt(randInterval)); err == nil {
		randomValue = n.String()
	} else {
		return fmt.Errorf(errFormatString, MetricNameRandomValue, err)
	}

	if m, err := model.NewMetricUnit(model.MetricTypeGauge, MetricNameRandomValue, randomValue); err == nil {
		*result = append(*result, m)
	} else {
		return fmt.Errorf(errFormatString, MetricNameRandomValue, err)
	}

	// PollCount
	if m, err := model.NewMetricUnit(model.MetricTypeCounter, MetricNamePollCount, "1"); err == nil {
		*result = append(*result, m)
	} else {
		return fmt.Errorf(errFormatString, MetricNamePollCount, err)
	}

	return nil
}

func collectOtherTypeMetrics(rtm *runtime.MemStats) (*[]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0)

	// GCCPUFraction
	fraction := strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)
	if m, err := model.NewMetricUnit(model.MetricTypeGauge, MetricNameGCCPUFraction, fraction); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, MetricNameGCCPUFraction, err)
	}

	// RandomValue
	var randInterval int64 = 100
	var randomValue string
	if n, err := rand.Int(rand.Reader, big.NewInt(randInterval)); err == nil {
		randomValue = n.String()
	} else {
		return nil, fmt.Errorf(errFormatString, MetricNameRandomValue, err)
	}

	if m, err := model.NewMetricUnit(model.MetricTypeGauge, MetricNameRandomValue, randomValue); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, MetricNameRandomValue, err)
	}

	// PollCount
	if m, err := model.NewMetricUnit(model.MetricTypeCounter, MetricNamePollCount, "1"); err == nil {
		result = append(result, m)
	} else {
		return nil, fmt.Errorf(errFormatString, MetricNamePollCount, err)
	}

	return &result, nil
}

// GetMemoryMetrics returns a list of the metrics.
func GetMemoryMetrics(resultChan chan *[]model.MetricUnit, errChan chan error) {
	const metricsCount = 3

	result := make([]model.MetricUnit, 0, metricsCount)
	var name string

	// CPUutilization1

	name = MetricNameCPUutiliz1

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
	name = MetricNameTotalMemory
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
	name = MetricNameFreeMemory
	if m, err := model.NewMetricUnit(model.MetricTypeGauge,
		name,
		strconv.FormatUint(virMem.Free, 10)); err == nil {
		result = append(result, m)
	} else {
		errChan <- fmt.Errorf(errFormatString, name, err)

		return
	}

	resultChan <- &result
}

// GetMemoryMetricsVariant returns a list of the metrics.
func GetMemoryMetricsVariant(resultChan chan *[]model.MetricUnit, errChan chan error) {
	const metricsCount = 3

	result := make([]model.MetricUnit, metricsCount)

	// CPUutilization1

	utilization, err := cpu.Percent(time.Duration(0), false)
	if err != nil {
		errChan <- fmt.Errorf(errFormatString, MetricNameCPUutiliz1, err)

		return
	}

	utilizationStValue := fmt.Sprintf("%v", utilization[0]) // strconv.FormatFloat(utilization[0], 'f', 2, 64)
	if mUnit, err := model.NewMetricUnit(model.MetricTypeGauge, MetricNameCPUutiliz1, utilizationStValue); err == nil {
		result[0] = mUnit
	} else {
		errChan <- fmt.Errorf(errFormatString, MetricNameCPUutiliz1, err)

		return
	}

	// TotalMemory
	virMem, err := mem.VirtualMemory()
	if err != nil {
		errChan <- fmt.Errorf(errFormatString, MetricNameTotalMemory, err)

		return
	}

	if mUnit, err := model.NewMetricUnit(model.MetricTypeGauge,
		MetricNameTotalMemory, strconv.FormatUint(virMem.Total, 10)); err == nil {
		result[1] = mUnit
	} else {
		errChan <- fmt.Errorf(errFormatString, MetricNameTotalMemory, err)

		return
	}

	// FreeMemory
	if mUnit, err := model.NewMetricUnit(model.MetricTypeGauge,
		MetricNameFreeMemory, strconv.FormatUint(virMem.Free, 10)); err == nil {
		result[2] = mUnit
	} else {
		errChan <- fmt.Errorf(errFormatString, MetricNameFreeMemory, err)

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
func getFieldValueVariant(mStat *runtime.MemStats, field string) string { //nolint:funlen
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
