package config

import (
	"testing"
)

/*
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
*/

func BenchmarkConfigStringVariantBuffer(b *testing.B) {
	cfg := ServerConfig{
		Config: Config{
			Address: "1",
			HashKey: "2",
		},
		StoreInterval:   3,
		FileStoragePath: "4",
		ConnectionDB:    "5",
		hasRestore:      true,
		Restore:         true,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a := cfg.StringVariantBuffer()
		_ = a
	}
}

func BenchmarkConfigString(b *testing.B) {
	cfg := ServerConfig{
		Config: Config{
			Address: "1",
			HashKey: "2",
		},
		StoreInterval:   3,
		FileStoragePath: "4",
		ConnectionDB:    "5",
		hasRestore:      true,
		Restore:         true,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a := cfg.String()
		_ = a
	}
}

func BenchmarkConfigStringCopy(b *testing.B) {
	cfg := ServerConfig{
		Config: Config{
			Address: "1",
			HashKey: "2",
		},
		StoreInterval:   3,
		FileStoragePath: "4",
		ConnectionDB:    "5",
		hasRestore:      true,
		Restore:         true,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a := cfg.StringVariantCopy()
		_ = a
	}
}
