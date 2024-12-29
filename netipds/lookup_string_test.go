package main_test

import (
	"net/netip"
	"testing"

	"github.com/aromatt/netipds"
	"local/iprbench/common"
)

var rt *netipds.PrefixMap[string]

func init() {
	rtb := netipds.PrefixMapBuilder[string]{}
	for i, route := range randomRoutes[:100_000] {
		rtb.Set(route, randomStrings[i])
	}
	rt = rtb.PrefixMap()
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
