package main_test

import (
	"local/iprbench/common"
	"testing"

	"github.com/aromatt/netipds"
)

func BenchmarkInsert(b *testing.B) {
	for k := 100; k <= 1_000_000; k *= 10 {
		rt := new(netipds.PrefixSetBuilder)
		for _, route := range tier1Routes[:k] {
			rt.Add(route)
		}

		name := "Insert into " + common.IntMap[k]
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				rt.Add(probe)
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	for k := 100; k <= 1_000_000; k *= 10 {
		rt := new(netipds.PrefixSetBuilder)
		for _, route := range tier1Routes[:k] {
			rt.Add(route)
		}

		name := "Delete from " + common.IntMap[k]
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				rt.Remove(probe)
			}
		})
	}
}

func BenchmarkInsertLazy(b *testing.B) {
	for k := 100; k <= 1_000_000; k *= 10 {
		rt := &netipds.PrefixSetBuilder{Lazy: true}
		for _, route := range tier1Routes[:k] {
			rt.Add(route)
		}

		name := "Insert into " + common.IntMap[k]
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				rt.Add(probe)
			}
		})
	}
}

func BenchmarkDeleteLazy(b *testing.B) {
	for k := 100; k <= 1_000_000; k *= 10 {
		rt := &netipds.PrefixSetBuilder{Lazy: true}
		for _, route := range tier1Routes[:k] {
			rt.Add(route)
		}

		name := "Delete from " + common.IntMap[k]
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				rt.Remove(probe)
			}
		})
	}
}
