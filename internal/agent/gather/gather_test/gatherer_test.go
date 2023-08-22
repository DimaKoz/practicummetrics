package gather_test_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
)

var testMetricsNameWantKeys = []string{
	gather.MetricNameAlloc,
	gather.MetricNameBuckHashSys,
	gather.MetricNameFrees,
	gather.MetricNameGCCPUFraction,
	gather.MetricNameGCSys,
	gather.MetricNameHeapAlloc,
	gather.MetricNameHeapIdle,
	gather.MetricNameHeapInuse,
	gather.MetricNameHeapObjects,
	gather.MetricNameHeapReleased,
	gather.MetricNameHeapSys,
	gather.MetricNameLastGC,
	gather.MetricNameLookups,
	gather.MetricNameMCacheInuse,
	gather.MetricNameMCacheSys,
	gather.MetricNameMSpanInuse,
	gather.MetricNameMSpanSys,
	gather.MetricNameMallocs,
	gather.MetricNameNextGC,
	gather.MetricNameNumForcedGC,
	gather.MetricNameNumGC,
	gather.MetricNameOtherSys,
	gather.MetricNamePauseTotalNs,
	gather.MetricNameStackInuse,
	gather.MetricNameStackSys,
	gather.MetricNameSys,
	gather.MetricNameTotalAlloc,
	gather.MetricNamePollCount,
	gather.MetricNameRandomValue,
}

var testMemoryMetricsNameWantKeys = []string{
	gather.MetricNameCPUutiliz1,
	gather.MetricNameTotalMemory,
	gather.MetricNameFreeMemory,
}

const (
	getMetrics              = 1
	getMetricsVariant       = 2
	getMemoryMetrics        = 3
	getMemoryMetricsVariant = 4
)

func startMetricsByState(t *testing.T, useState int, resultChan chan *[]model.MetricUnit, errChan chan error) {
	t.Helper()
	switch {
	case useState == getMetrics:
		go gather.GetMetrics(resultChan, errChan)
	case useState == getMetricsVariant:
		go gather.GetMetricsVariant(resultChan, errChan)
	case useState == getMemoryMetrics:
		go gather.GetMemoryMetrics(resultChan, errChan)
	case useState == getMemoryMetricsVariant:
		go gather.GetMemoryMetricsVariant(resultChan, errChan)
	default:
		require.Fail(t, "Unknown test.use")
	}
}

func TestGetMetrics(t *testing.T) {
	tests := []struct {
		name     string
		wantKeys []string
		use      int
	}{
		{
			name: "test keys", wantKeys: testMetricsNameWantKeys, use: getMetrics,
		},
		{
			name: "test keys variant", wantKeys: testMetricsNameWantKeys, use: getMetricsVariant,
		},
		{
			name: "test memory keys", wantKeys: testMemoryMetricsNameWantKeys, use: getMemoryMetrics,
		},
		{
			name: "test memory keys variant", wantKeys: testMemoryMetricsNameWantKeys, use: getMemoryMetricsVariant,
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			metricsCh := make(chan *[]model.MetricUnit)
			errCh := make(chan error)
			startMetricsByState(t, test.use, metricsCh, errCh)

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
						t.Errorf("metrics = %v, want %v", got, test.wantKeys)
						checkMetricsName(t, test.wantKeys, got)
					} else {
						checkMetricsName(t, test.wantKeys, got)
					}

					break ForLoop
				}
			}
			close(metricsCh)
			close(errCh)
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
			assert.Fail(testing, "GetMetrics() - no expected metric", "we want %s but absentee", wantKey)
		}
	}
}

/*

goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/gather/gather_test
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkGetMetrics
BenchmarkGetMetrics-8          	   26454	     43101 ns/op	   12646 B/op	      69 allocs/op
BenchmarkGetMetricsVariant
BenchmarkGetMetricsVariant-8   	   42862	     27827 ns/op	    2625 B/op	      37 allocs/op

*/

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
				assert.NotNil(b, got)

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
				assert.NotNil(b, got)

				break ForLoop
			}
		}
		close(metricsCh)
		close(errCh)
	}
}

func BenchmarkGetMemoryMetrics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metricsCh := make(chan *[]model.MetricUnit)
		errCh := make(chan error)
		go gather.GetMemoryMetrics(metricsCh, errCh)

	ForLoop:
		for {
			select {
			case err := <-errCh:
				b.Error(err)

				break ForLoop
			case got := <-metricsCh:
				assert.NotNil(b, got)

				break ForLoop
			}
		}
		close(metricsCh)
		close(errCh)
	}
}

func BenchmarkGetMemoryMetricsVariant(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metricsCh := make(chan *[]model.MetricUnit)
		errCh := make(chan error)
		go gather.GetMemoryMetricsVariant(metricsCh, errCh)

	ForLoop:
		for {
			select {
			case err := <-errCh:
				b.Error(err)

				break ForLoop
			case got := <-metricsCh:
				assert.NotNil(b, got)

				break ForLoop
			}
		}
		close(metricsCh)
		close(errCh)
	}
}
