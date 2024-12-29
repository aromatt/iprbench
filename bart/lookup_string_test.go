package main_test

import (
	"net/netip"
	"testing"

	"github.com/gaissmai/bart"
	"local/iprbench/common"
)

var rt = new(bart.Table[any])

func init() {
	for i, route := range randomRoutes[:100_000] {
		rt.Insert(route, randomStrings[i])
	}
}

func BenchmarkLpmRandomPfxsStr100_000(b *testing.B) {
	benchmarks := []struct {
		name    string
		routes  []netip.Prefix
		strings []string
		fn      func([]netip.Prefix) netip.Addr
	}{
		{"RandomMatchIP4", randomRoutes[:100_000], randomStrings, common.MatchIP4},
		{"RandomMatchIP6", randomRoutes[:100_000], randomStrings, common.MatchIP6},
		{"RandomMissIP4", randomRoutes[:100_000], randomStrings, common.MissIP4},
		{"RandomMissIP6", randomRoutes[:100_000], randomStrings, common.MissIP6},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ip := bm.fn(bm.routes)
			pfx := netip.PrefixFrom(ip, ip.BitLen())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, sink = rt.Get(pfx)
			}
		})
	}
}
