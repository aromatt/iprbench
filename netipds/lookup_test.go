package main_test

import (
	"net/netip"
	"testing"

	"local/iprbench/common"

	"github.com/aromatt/netipds"
)

var rt1b = new(netipds.PrefixSetBuilder)
var rt2b = new(netipds.PrefixSetBuilder)

func init() {
	for _, route := range tier1Routes {
		rt1b.Add(route)
	}
}

func init() {
	for _, route := range randomRoutes[:100_000] {
		rt2b.Add(route)
	}
}

func BenchmarkLpmTier1Pfxs(b *testing.B) {

	benchmarks := []struct {
		name   string
		routes []netip.Prefix
		fn     func([]netip.Prefix) netip.Addr
	}{
		{"RandomMatchIP4", tier1Routes, common.MatchIP4},
		{"RandomMatchIP6", tier1Routes, common.MatchIP6},
		{"RandomMissIP4", tier1Routes, common.MissIP4},
		{"RandomMissIP6", tier1Routes, common.MissIP6},
	}

	rt1 := rt1b.PrefixSet()
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ip := bm.fn(bm.routes)
			pfx := netip.PrefixFrom(ip, ip.BitLen())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = rt1.Contains(pfx)
			}
		})
	}
}

func BenchmarkLpmRandomPfxs100_000(b *testing.B) {

	benchmarks := []struct {
		name   string
		routes []netip.Prefix
		fn     func([]netip.Prefix) netip.Addr
	}{
		{"RandomMatchIP4", randomRoutes[:100_000], common.MatchIP4},
		{"RandomMatchIP6", randomRoutes[:100_000], common.MatchIP6},
		{"RandomMissIP4", randomRoutes[:100_000], common.MissIP4},
		{"RandomMissIP6", randomRoutes[:100_000], common.MissIP6},
	}

	rt2 := rt2b.PrefixSet()
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ip := bm.fn(bm.routes)
			pfx := netip.PrefixFrom(ip, ip.BitLen())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = rt2.Contains(pfx)
			}
		})
	}
}