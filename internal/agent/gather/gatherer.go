package gather

import (
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

// GetMetrics returns a list of the metrics
func GetMetrics() *[]model.MetricUnit {
	result := make([]model.MetricUnit, 0)

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	var m *model.MetricUnit
	var err *error2.RequestError

	// Alloc
	alloc := strconv.FormatUint(rtm.Alloc, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "Alloc", alloc); err == nil && m != nil {
		result = append(result, *m)
	}

	// BuckHashSys
	buckHashSys := strconv.FormatUint(rtm.BuckHashSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "BuckHashSys", buckHashSys); err == nil && m != nil {
		result = append(result, *m)
	}

	// Frees
	frees := strconv.FormatUint(rtm.Frees, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "Frees", frees); err == nil && m != nil {
		result = append(result, *m)
	}

	// GCCPUFraction
	fraction := strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "GCCPUFraction", fraction); err == nil && m != nil {
		result = append(result, *m)
	}

	// GCSys
	gcsys := strconv.FormatUint(rtm.GCSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "GCSys", gcsys); err == nil && m != nil {
		result = append(result, *m)
	}

	// HeapAlloc
	heapAlloc := strconv.FormatUint(rtm.HeapAlloc, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "HeapAlloc", heapAlloc); err == nil && m != nil {
		result = append(result, *m)
	}

	// HeapIdle
	heapIdle := strconv.FormatUint(rtm.HeapAlloc, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "HeapIdle", heapIdle); err == nil && m != nil {
		result = append(result, *m)
	}

	// HeapInuse
	heapInuse := strconv.FormatUint(rtm.HeapInuse, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "HeapInuse", heapInuse); err == nil && m != nil {
		result = append(result, *m)
	}

	// HeapObjects
	heapObjects := strconv.FormatUint(rtm.HeapObjects, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "HeapObjects", heapObjects); err == nil && m != nil {
		result = append(result, *m)
	}

	// HeapReleased
	heapReleased := strconv.FormatUint(rtm.HeapReleased, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "HeapReleased", heapReleased); err == nil && m != nil {
		result = append(result, *m)
	}

	// HeapSys
	heapSys := strconv.FormatUint(rtm.HeapSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "HeapSys", heapSys); err == nil && m != nil {
		result = append(result, *m)
	}

	// LastGC
	lastGC := strconv.FormatUint(rtm.LastGC/1_000_000, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "LastGC", lastGC); err == nil && m != nil {
		result = append(result, *m)
	}

	// Lookups
	lookups := strconv.FormatUint(rtm.Lookups, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "Lookups", lookups); err == nil && m != nil {
		result = append(result, *m)
	}

	// MCacheInuse
	mCacheInuse := strconv.FormatUint(rtm.MCacheInuse, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "MCacheInuse", mCacheInuse); err == nil && m != nil {
		result = append(result, *m)
	}

	// MCacheSys
	mCacheSys := strconv.FormatUint(rtm.MCacheSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "MCacheSys", mCacheSys); err == nil && m != nil {
		result = append(result, *m)
	}

	// MSpanInuse
	mSpanInuse := strconv.FormatUint(rtm.MSpanInuse, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "MSpanInuse", mSpanInuse); err == nil && m != nil {
		result = append(result, *m)
	}

	// MSpanSys
	mSpanSys := strconv.FormatUint(rtm.MSpanSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "MSpanSys", mSpanSys); err == nil && m != nil {
		result = append(result, *m)
	}

	// Mallocs
	malloc := strconv.FormatUint(rtm.Mallocs, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "Mallocs", malloc); err == nil && m != nil {
		result = append(result, *m)
	}

	// NextGC
	nextGC := strconv.FormatUint(rtm.NextGC, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "NextGC", nextGC); err == nil && m != nil {
		result = append(result, *m)
	}

	// NumForcedGC
	numForcedGC := strconv.FormatUint(uint64(rtm.NumForcedGC), 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "NumForcedGC", numForcedGC); err == nil && m != nil {
		result = append(result, *m)
	}

	// NumGC
	numGC := strconv.FormatUint(uint64(rtm.NumGC), 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "NumGC", numGC); err == nil && m != nil {
		result = append(result, *m)
	}

	// OtherSys
	otherSys := strconv.FormatUint(rtm.OtherSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "OtherSys", otherSys); err == nil && m != nil {
		result = append(result, *m)
	}

	// PauseTotalNs
	pauseTotalNs := strconv.FormatUint(rtm.PauseTotalNs, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "PauseTotalNs", pauseTotalNs); err == nil && m != nil {
		result = append(result, *m)
	}

	// StackInuse
	stackInuse := strconv.FormatUint(rtm.StackInuse, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "StackInuse", stackInuse); err == nil && m != nil {
		result = append(result, *m)
	}

	// StackSys
	stackSys := strconv.FormatUint(rtm.StackSys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "StackSys", stackSys); err == nil && m != nil {
		result = append(result, *m)
	}

	// Sys
	sys := strconv.FormatUint(rtm.Sys, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "Sys", sys); err == nil && m != nil {
		result = append(result, *m)
	}

	// TotalAlloc
	totalAlloc := strconv.FormatUint(rtm.TotalAlloc, 10)
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "TotalAlloc", totalAlloc); err == nil && m != nil {
		result = append(result, *m)
	}

	// RandomValue
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomValue := strconv.Itoa(r1.Intn(100))
	if m, err = model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", randomValue); err == nil && m != nil {
		result = append(result, *m)
	}

	// PollCount
	if m, err = model.NewMetricUnit(model.MetricTypeCounter, "PollCount", "1"); err == nil && m != nil {
		result = append(result, *m)
	}

	return &result
}
