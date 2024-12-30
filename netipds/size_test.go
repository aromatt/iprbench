package main_test

import (
	"runtime"
	"testing"

	"local/iprbench/common"

	"github.com/aromatt/netipds"
)

func BenchmarkTier1PfxSize(b *testing.B) {
	var startMem, endMem runtime.MemStats

	for k := 100; k <= 1_000_000; k *= 10 {
		tree := new(netipds.PrefixSetBuilder)
		runtime.GC()
		runtime.ReadMemStats(&startMem)

		b.Run(common.IntMap[k], func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for _, cidr := range tier1Routes[:k] {
					tree.Add(cidr)
				}
			}
			runtime.GC()
			runtime.ReadMemStats(&endMem)

			b.ReportMetric(float64(endMem.HeapAlloc-startMem.HeapAlloc), "bytes")
			b.ReportMetric(0, "ns/op")
		})
	}
}

func BenchmarkRandomPfx4Size(b *testing.B) {
	var startMem, endMem runtime.MemStats

	for k := 100; k <= 1_000_000; k *= 10 {
		tree := new(netipds.PrefixSetBuilder)
		runtime.GC()
		runtime.ReadMemStats(&startMem)

		b.Run(common.IntMap[k], func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for _, cidr := range randomRoutes4[:k] {
					tree.Add(cidr)
				}
			}
			runtime.GC()
			runtime.ReadMemStats(&endMem)

			b.ReportMetric(float64(endMem.HeapAlloc-startMem.HeapAlloc), "bytes")
			b.ReportMetric(0, "ns/op")
		})
	}
}

func BenchmarkRandomPfx6Size(b *testing.B) {
	var startMem, endMem runtime.MemStats

	for k := 100; k <= 1_000_000; k *= 10 {
		tree := new(netipds.PrefixSetBuilder)
		runtime.GC()
		runtime.ReadMemStats(&startMem)

		b.Run(common.IntMap[k], func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for _, cidr := range randomRoutes6[:k] {
					tree.Add(cidr)
				}
			}
			runtime.GC()
			runtime.ReadMemStats(&endMem)

			b.ReportMetric(float64(endMem.HeapAlloc-startMem.HeapAlloc), "bytes")
			b.ReportMetric(0, "ns/op")
		})
	}
}

func BenchmarkRandomPfxSize(b *testing.B) {
	var startMem, endMem runtime.MemStats

	for k := 100; k <= 1_000_000; k *= 10 {
		tree := new(netipds.PrefixSetBuilder)
		runtime.GC()
		runtime.ReadMemStats(&startMem)

		b.Run(common.IntMap[k], func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for _, cidr := range randomRoutes[:k] {
					tree.Add(cidr)
				}
			}
			runtime.GC()
			runtime.ReadMemStats(&endMem)

			b.ReportMetric(float64(endMem.HeapAlloc-startMem.HeapAlloc), "bytes")
			b.ReportMetric(0, "ns/op")
		})
	}
}
