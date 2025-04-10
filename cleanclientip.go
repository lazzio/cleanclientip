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

	// X-Forwarded-For can contain multiple IPs, we only want the first one
	ips := strings.Split(xff, ",")
	if len(ips) == 0 {
		a.next.ServeHTTP(rw, req)
		return
	}

	// Get each IP address and remove ports
	for i, ip := range ips {
		ip = strings.TrimSpace(ip)
		// Remove port if present
		ips[i] = strings.Split(ip, ":")[0]
	}

	// Set the first IP address as the remote address
	req.RemoteAddr = ips[0]

	// Set the X-Forwarded-For header to the cleaned IPs
	req.Header.Set("X-Forwarded-For", strings.Join(ips, ", "))

	// Set the X-Real-Ip header to the first IP address
	req.Header.Set("X-Real-Ip", ips[0])

	// Call the next handler
	a.next.ServeHTTP(rw, req)
}
