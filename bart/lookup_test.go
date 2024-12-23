package main_test

import (
	"net/netip"
	"testing"

	"local/iprbench/common"

	"github.com/gaissmai/bart"
	profile "github.com/pkg/profile"
)

var rt1 = new(bart.Table[any])
var rt2 = new(bart.Table[any])

func init() {
	for _, route := range tier1Routes {
		rt1.Insert(route, nil)
	}
}

func init() {
	for _, route := range randomRoutes[:100_000] {
		rt2.Insert(route, nil)
	}
}

func BenchmarkLpmTier1Pfxs(b *testing.B) {

	defer profile.Start(profile.ProfilePath(".")).Stop()
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

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ip := bm.fn(bm.routes)
			b.ResetTimer()
			for range b.N {
				_, sink = rt1.Lookup(ip)
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

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ip := bm.fn(bm.routes)
			b.ResetTimer()
			for range b.N {
				_, sink = rt2.Lookup(ip)
			}
		})
	}
}
