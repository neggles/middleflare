package cfaddrs

import (
	"log"
	"net/netip"
)

func CheckSourceAddr(ip netip.Addr, prefixes ...netip.Prefix) bool {
	if len(prefixes) == 0 {
		prefixes = FallbackAddresses()
	}
	for _, network := range prefixes {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func ParsePrefixes(cidrs []string) []netip.Prefix {
	prefixes := make([]netip.Prefix, len(cidrs))
	for i, cidr := range cidrs {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			log.Fatal(err)
		}
		prefixes[i] = prefix
	}
	return prefixes[:]
}
