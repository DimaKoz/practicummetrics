package gather

import (
	"runtime"
	"testing"
)

/*

goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/gather
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkGetFieldValueUint64
BenchmarkGetFieldValueUint64-8   	  203800	      5629 ns/op	     408 B/op	      49 allocs/op
BenchmarkGetFieldValue
BenchmarkGetFieldValue-8         	 1000000	      1005 ns/op	     200 B/op	      23 allocs/op

BenchmarkGetFieldValueUint64 - see 'deadcode_grave' branch
*/

func BenchmarkGetFieldValue(b *testing.B) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	rtmPtr := &rtm
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, name := range metricsName {
			getFieldValueVariant(rtmPtr, name)
		}
	}
}
