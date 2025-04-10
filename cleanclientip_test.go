package cleanclientip_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lazzio/cleanclientip"
)

func TestCleanClientIp(t *testing.T) {
	cfg := cleanclientip.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := cleanclientip.New(ctx, next, cfg, "cleanclientip")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test with multiple IPs in X-Forwarded-For header
	req.Header.Set("X-Forwarded-For", "10.0.0.1:1234, 10.0.0.2, 10.0.0.3:5678, 10.0.0.4")

	handler.ServeHTTP(recorder, req)
	assertHeader(t, req, "X-Forwarded-For", "10.0.0.1, 10.0.0.2, 10.0.0.3, 10.0.0.4")
	assertHeader(t, req, "X-Real-Ip", "10.0.0.1")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
