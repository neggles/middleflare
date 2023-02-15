package cfaddrs // import "github.com/neggles/middleflare/cfaddrs"

import (
	"log"
	"net/netip"
)

// CheckSourceAddr checks if the given IP address is in the given list of prefixes.
func CheckSourceAddr(ip netip.Addr, prefixes ...netip.Prefix) bool {
	if len(prefixes) == 0 {
		prefixes = CloudflareAddresses()
	}
	for _, network := range prefixes {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// ParsePrefixes parses a list of CIDR strings into a list of netip.Prefix.
func ParsePrefixes(cidrs []string) []netip.Prefix {
	prefixes := make([]netip.Prefix, len(cidrs))
	for i, cidr := range cidrs {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			log.Fatal(err)
		}
		prefixes[i] = prefix
	}
	return prefixes
}
