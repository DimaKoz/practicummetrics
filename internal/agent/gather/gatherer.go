package gather

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Constants to store metrics names.
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

// metricsName contains a collection of metrics.
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

// errFormatString a template of an error.
const errFormatString = "error while collecting metrics with: \n can't get '%s' metric by %w "

// collectUintMetricsVariant collects model.MetricUnit from runtime.MemStats.
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

// collectOtherTypeMetricsVariant collects model.MetricUnit from runtime.MemStats.
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
	// utilization[0]
	// 6.444818871322972

	// fmt.Sprintf is better than strconv.FormatFloat(utilization[0], 'f', 2, 64)
	// see BenchmarkUtilizationConvertValue
	utilizationStValue := fmt.Sprintf("%v", utilization[0])
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

// getFieldValueVariant gets field value from runtime.MemStats.
func getFieldValueVariant(mStat *runtime.MemStats, field string) string { //nolint:funlen,cyclop
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
