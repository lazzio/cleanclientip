package cleanclientip

import (
	"context"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type CleanClientIp struct {
	next http.Handler
	name string
}

// New created a new CleanClientIp plugin.
func New(ctx context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &CleanClientIp{
		next: next,
		name: name,
	}, nil
}

func (a *CleanClientIp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	xff := req.Header.Get("X-Forwarded-For")
	if xff == "" {
		a.next.ServeHTTP(rw, req)
		return
	}

	ips := strings.Split(xff, ",")
	if len(ips) == 0 {
		a.next.ServeHTTP(rw, req)
		return
	}

	// Clean IPs by removing spaces and ports
	// and store them in cleanedIPs
	cleanedIPs := make([]string, len(ips))
	for i, ip := range ips {
		ip = strings.TrimSpace(ip)
		// Enlever le port s'il est présent
		cleanedIPs[i] = strings.Split(ip, ":")[0]
	}

	// Update X-Forwarded-For header with cleaned IPs
	req.Header.Set("X-Forwarded-For", strings.Join(cleanedIPs, ", "))

	// Set X-Real-Ip header with the first cleaned IP
	req.Header.Set("X-Real-Ip", cleanedIPs[0])

	// Call the next handler
	// Pass the modified request to the next handler
	// This is important to ensure that the changes are reflected in the request
	// that is passed to the next handler
	// and that the original request is not modified.
	a.next.ServeHTTP(rw, req)
}
