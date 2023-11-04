package gather

import (
	"fmt"
	"runtime"
	"strconv"
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

/*
goos: linux
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/gather
cpu: 11th Gen Intel(R) Core(TM) i5-11400 @ 2.60GHz
BenchmarkUtilizationConvertValue
BenchmarkUtilizationConvertValue/Sprintf()
BenchmarkUtilizationConvertValue/Sprintf()-12         	 7472413	       157.1 ns/op	      32 B/op	       2 allocs/op
BenchmarkUtilizationConvertValue/FormatFloat()
BenchmarkUtilizationConvertValue/FormatFloat()-12     	 8379586	       142.1 ns/op	      28 B/op	       2 allocs/op
*/

func BenchmarkUtilizationConvertValue(b *testing.B) {
	utilization := 6.444818871322972
	b.ResetTimer()
	b.Run("Sprintf()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := fmt.Sprintf("%v", utilization)
			_ = a
		}
	})
	b.Run("FormatFloat()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := strconv.FormatFloat(utilization, 'f', 2, 64)
			_ = a
		}
	})
}
