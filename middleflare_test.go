package middleflare_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neggles/middleflare"
)

func TestDemo(t *testing.T) {
	cfg := middleflare.CreateConfig()
	cfg.TrustedProxies = []string{}
	cfg.IncludeDefault = true

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := middleflare.New(ctx, next, cfg, "middleflare")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	req.Header.Add("CF-Connecting-IP", "1.2.3.4")
	req.Header.Add("CF-Visitor", `{"scheme":"https"}`)
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-For", "1.2.3.4 4.3.2.1 8.8.8.8 7.7.7.7")
	req.RemoteAddr = "127.0.0.1:52342"

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)
	assertHeader(t, req, "X-Forwarded-For", "1.2.3.4")
	assertHeader(t, req, "X-Forwarded-Proto", "https")
	assertHeader(t, req, "X-Trusted-Proxy", "127.0.0.1")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
