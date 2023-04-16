package gather

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"testing"
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
				"RandomValue"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := GetMetrics(); len(got) != len(tt.wantKeys) {

				t.Errorf("GetMetrics() = %v, want %v", got, tt.wantKeys)
				checkMetricsName(t, tt.wantKeys, &got)
			} else {
				checkMetricsName(t, tt.wantKeys, &got)
			}
		})
	}
}

func checkMetricsName(t *testing.T, wantKeys []string, got *[]model.MetricUnit) {
	for _, k := range wantKeys {
		isPresent := false
		for _, kk := range *got {
			if kk.Name == k {
				isPresent = true
				break
			}
		}
		if !isPresent {
			t.Errorf("GetMetrics() -  we want %v but absentee", k)
		}
	}

}
