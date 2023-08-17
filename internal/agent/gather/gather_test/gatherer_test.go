package gather_test_test

import (
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/stretchr/testify/assert"
)

var testMetricsNameWantKeys = []string{
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
}

func TestGetMetrics(t *testing.T) {
	tests := []struct {
		name      string
		wantKeys  []string
		isVariant bool
	}{
		{
			name:      "test keys",
			wantKeys:  testMetricsNameWantKeys,
			isVariant: false,
		},
		{
			name:      "test keys variant",
			wantKeys:  testMetricsNameWantKeys,
			isVariant: true,
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			metricsCh := make(chan *[]model.MetricUnit)
			errCh := make(chan error)
			if test.isVariant {
				go gather.GetMetricsVariant(metricsCh, errCh)
			} else {
				go gather.GetMetrics(metricsCh, errCh)
			}

		ForLoop:
			for {
				select {
				case err := <-errCh:
					assert.NoError(t, err)

					break ForLoop
				case got := <-metricsCh:
					assert.NotNil(t, got)
					if got != nil && len(*got) != len(test.wantKeys) {
						assert.Equal(t, len(test.wantKeys), len(*got))
						t.Errorf("GetMetrics() = %v, want %v", got, test.wantKeys)
						checkMetricsName(t, test.wantKeys, got)
					} else {
						checkMetricsName(t, test.wantKeys, got)
					}

					break ForLoop
				}
			}
		})
	}
}

func checkMetricsName(testing assert.TestingT, wantKeys []string, got *[]model.MetricUnit) {
	for _, wantKey := range wantKeys {
		isPresent := false

		for _, kk := range *got {
			if kk.Name == wantKey {
				isPresent = true

				break
			}
		}

		if !isPresent {
			assert.Fail(testing, "GetMetrics() -  we want %v but absentee", wantKey)
		}
	}
}

func BenchmarkGetMetrics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metricsCh := make(chan *[]model.MetricUnit)
		errCh := make(chan error)
		go gather.GetMetrics(metricsCh, errCh)

	ForLoop:
		for {
			select {
			case err := <-errCh:
				b.Error(err)

				break ForLoop
			case got := <-metricsCh:
				if got == nil {
					b.Fatalf("got == ni")
				}
				if got != nil && len(*got) != len(testMetricsNameWantKeys) {
					assert.Equal(b, len(testMetricsNameWantKeys), len(*got))
					b.Errorf("GetMetrics() = %v, want %v", got, testMetricsNameWantKeys)
					checkMetricsName(b, testMetricsNameWantKeys, got)
				} else {
					checkMetricsName(b, testMetricsNameWantKeys, got)
				}

				break ForLoop
			}
		}
		close(metricsCh)
		close(errCh)
	}
}

func BenchmarkGetMetricsVariant(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metricsCh := make(chan *[]model.MetricUnit)
		errCh := make(chan error)
		go gather.GetMetricsVariant(metricsCh, errCh)

	ForLoop:
		for {
			select {
			case err := <-errCh:
				b.Error(err)

				break ForLoop
			case got := <-metricsCh:
				if got == nil {
					b.Fatalf("got == ni")
				}
				if got != nil && len(*got) != len(testMetricsNameWantKeys) {
					assert.Equal(b, len(testMetricsNameWantKeys), len(*got))
					b.Errorf("GetMetricsVariant() = %v, want %v", got, testMetricsNameWantKeys)
					checkMetricsName(b, testMetricsNameWantKeys, got)
				} else {
					checkMetricsName(b, testMetricsNameWantKeys, got)
				}

				break ForLoop
			}
		}
		close(metricsCh)
		close(errCh)
	}
}
