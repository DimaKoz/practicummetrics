# Бенчмарки
<pre>
goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/common/repository
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkLoad
BenchmarkLoad/Load()
BenchmarkLoad/Load()-8         	   52819	     23062 ns/op	    1512 B/op	      21 allocs/op
BenchmarkLoad/LoadVariant()
BenchmarkLoad/LoadVariant()-8  	   71215	     15445 ns/op	    1384 B/op	      12 allocs/op
BenchmarkSave
BenchmarkSave/Save()
BenchmarkSave/Save()-8         	    4951	    224940 ns/op	     232 B/op	       5 allocs/op
BenchmarkSave/SaveVariant()
BenchmarkSave/SaveVariant()-8  	    5010	    232807 ns/op	     226 B/op	       5 allocs/op
</pre>

<pre>
goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/common/config
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkConfigStringVariantBuffer
BenchmarkConfigStringVariantBuffer-8   	 7398451	       161.9 ns/op	     192 B/op	       2 allocs/op
BenchmarkConfigString
BenchmarkConfigString-8                	 2832891	       426.9 ns/op	     160 B/op	       5 allocs/op
BenchmarkConfigStringCopy
BenchmarkConfigStringCopy-8            	10617492	       120.4 ns/op	     192 B/op	       2 allocs/op
</pre>

<pre>
goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/gather/gather_test
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkGetMetrics
BenchmarkGetMetrics-8          	   26454	     43101 ns/op	   12646 B/op	      69 allocs/op
BenchmarkGetMetricsVariant
BenchmarkGetMetricsVariant-8   	   42862	     27827 ns/op	    2625 B/op	      37 allocs/op
</pre>

<pre>
goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/gather
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkGetFieldValueUint64
BenchmarkGetFieldValueUint64-8   	  203800	      5629 ns/op	     408 B/op	      49 allocs/op
BenchmarkGetFieldValue
BenchmarkGetFieldValue-8         	 1000000	      1005 ns/op	     200 B/op	      23 allocs/op
</pre>

<pre>
goos: darwin
goarch: amd64
pkg: github.com/DimaKoz/practicummetrics/internal/agent/sender
cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
BenchmarkAppendHash
BenchmarkAppendHash-8        	  608454	      2139 ns/op	     728 B/op	      12 allocs/op
BenchmarkAppendHashGoccy
BenchmarkAppendHashGoccy-8   	  709033	      1666 ns/op	     680 B/op	      11 allocs/op
</pre>