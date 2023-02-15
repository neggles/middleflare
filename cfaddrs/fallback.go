package cfaddrs

import (
	"net/netip"
)

// CloudflareIPv4 returns a list of IPv4 prefixes that are used by Cloudflare.
func CloudflareIPv4() []netip.Prefix {
	fallbackIPv4 := [...]string{
		"127.0.0.1/32", // include localhost for testing
		"10.16.0.0/20",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
	}
	return ParsePrefixes(fallbackIPv4[:])
}

// CloudflareIPv6 returns a list of IPv6 prefixes that are used by Cloudflare.
func CloudflareIPv6() []netip.Prefix {
	fallbackIPv6 := [...]string{
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}
	return ParsePrefixes(fallbackIPv6[:])
}

// CloudflareAddresses returns a list of both IPv4 and IPv6 prefixes that are used by Cloudflare.
func CloudflareAddresses() []netip.Prefix {
	return append(CloudflareIPv4(), CloudflareIPv6()...)
}
