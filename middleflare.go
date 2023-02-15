package middleflare

import (
	"context"
	"net/http"
	"net/netip"
	"strings"

	"github.com/neggles/middleflare/cfaddrs"
)

const (
	XRealIP         = "X-Real-IP"
	XForwardedFor   = "X-Forwarded-For"
	XForwardedProto = "X-Forwarded-Proto"
	XForwardedHost  = "X-Forwarded-Host"
	XTrustedProxy   = "X-Trusted-Proxy"
	CFConnectingIP  = "CF-Connecting-IP"
	CFVisitor       = "CF-Visitor"
)

// Config the plugin configuration.
type Config struct {
	TrustedProxies []string `json:"trustedProxies,omitempty"`
	IncludeDefault bool     `json:"includeDefault,omitempty"`
}

type CheckResult struct {
	IsValid   bool
	IsTrusted bool
	ProxyAddr netip.Addr
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		TrustedProxies: []string{},
		IncludeDefault: true,
	}
}

type CloudflareHeaders struct {
}

// CFHeaderWriter is a plugin that maps CF-Connecting-IP to X-Real-IP and X-Forwarded-For.
type CFHeaderWriter struct {
	next          http.Handler
	name          string
	trustPrefixes []netip.Prefix
}

// New created a new CFHeaderWriter plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config == nil {
		config = CreateConfig()
	}

	var trustPrefixes []netip.Prefix

	// If IncludeDefault is true, then we add the default fallback addresses.
	if config.IncludeDefault {
		trustPrefixes = cfaddrs.FallbackAddresses()
	}

	// If TrustedProxies is not empty, then we add the user defined addresses.
	if len(config.TrustedProxies) > 0 {
		trustPrefixes = append(trustPrefixes, cfaddrs.ParsePrefixes(config.TrustedProxies)...)
	}

	// If we have no addresses to trust, then we return an error.
	return &CFHeaderWriter{
		next:          next,
		name:          name,
		trustPrefixes: trustPrefixes,
	}, nil
}

func (writer *CFHeaderWriter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Check if the request is coming from a trusted proxy.
	checkResult := writer.checkSourceAddr(req.RemoteAddr)

	// If the remote address is invalid, then we return an error.
	if !checkResult.IsValid {
		http.Error(rw, "Invalid remote address", http.StatusInternalServerError)
		return
	}

	// Set headers if the request is coming from a trusted proxy.
	if checkResult.IsTrusted {
		// Set X-Trusted-Proxy header if we have a valid proxy address.
		if checkResult.ProxyAddr.IsValid() {
			req.Header.Set(XTrustedProxy, checkResult.ProxyAddr.String())
		}

		// Set X-Real-IP and X-Forwarded-For headers if the CF-Connecting-IP header is set.
		if req.Header.Get(CFConnectingIP) != "" {
			req.Header.Set(XRealIP, req.Header.Get(CFConnectingIP))
			req.Header.Set(XForwardedFor, req.Header.Get(CFConnectingIP))
		}
	}

	// Pass the request to the next middleware.
	writer.next.ServeHTTP(rw, req)
}

func (writer *CFHeaderWriter) checkSourceAddr(remoteAddr string) *CheckResult {
	// Split the remote address into the IP and port, and then take the IP.
	strIP := strings.Split(remoteAddr, ":")[0]

	// Parse the IP address.
	addr, err := netip.ParseAddr(strIP)
	if err != nil {
		return &CheckResult{
			IsValid:   false,
			IsTrusted: false,
		}
	}

	// Check if the address is in the trusted proxy list.
	if len(writer.trustPrefixes) > 0 {
		for _, network := range writer.trustPrefixes {
			if network.Contains(addr) {
				return &CheckResult{
					IsValid:   true,
					IsTrusted: true,
					ProxyAddr: addr,
				}
			}
		}
	}

	// If we get here, then the remote address is not trusted, or we trust no proxies.
	return &CheckResult{
		IsValid:   true,
		IsTrusted: false,
		ProxyAddr: addr,
	}
}
