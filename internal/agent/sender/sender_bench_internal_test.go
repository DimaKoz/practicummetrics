package sender

import (
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
	goccyj "github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

/*

goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/sender
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkAppendHash
BenchmarkAppendHash-8        	  608454	      2139 ns/op	     728 B/op	      12 allocs/op
BenchmarkAppendHashGoccy
BenchmarkAppendHashGoccy-8   	  709033	      1666 ns/op	     680 B/op	      11 allocs/op

*/

func BenchmarkAppendHash(b *testing.B) {
	metUnit, _ := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", "4321")
	emptyMetrics := model.NewEmptyMetrics()
	emptyMetrics.UpdateByMetricUnit(metUnit)
	request := resty.New().R()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := appendHash(request, "12345", emptyMetrics)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkAppendHashGoccy(b *testing.B) {
	metUnit, _ := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", "4321")
	emptyMetrics := model.NewEmptyMetrics()
	emptyMetrics.UpdateByMetricUnit(metUnit)
	request := resty.New().R()
	body, err := goccyj.Marshal(emptyMetrics)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		appendHashOtherMarshaling(request, "12345", body)
	}
}
