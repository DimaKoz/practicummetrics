package gather_test_test

import (
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
)

func TestGetMetrics(t *testing.T) {
	tests := []struct {
		name     string
		wantKeys []string
	}{
		{
			name: "test keys",
			wantKeys: []string{
				"Alloc",
				"BuckHashSys",
				"Frees",
				"GCCPUFraction",
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
				"NumForcedGC",
				"NumGC",
				"OtherSys",
				"PauseTotalNs",
				"StackInuse",
				"StackSys",
				"Sys",
				"TotalAlloc",
				"PollCount",
				"RandomValue",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, _ := gather.GetMetrics(); got != nil && len(*got) != len(test.wantKeys) {
				t.Errorf("GetMetrics() = %v, want %v", got, test.wantKeys)
				checkMetricsName(t, test.wantKeys, got)
			} else {
				checkMetricsName(t, test.wantKeys, got)
			}
		})
	}
}

func checkMetricsName(t *testing.T, wantKeys []string, got *[]model.MetricUnit) {
	t.Helper()

	for _, wantKey := range wantKeys {
		isPresent := false

		for _, kk := range *got {
			if kk.Name == wantKey {
				isPresent = true

				break
			}
		}

		if !isPresent {
			t.Errorf("GetMetrics() -  we want %v but absentee", wantKey)
		}
	}
}
