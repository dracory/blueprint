package httpsredirect

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/dracory/rtr"
)

var devEnvironments = map[string]bool{
	"development": true,
	"local":       true,
	"testing":     true,
}

// isDevEnvironment checks if the current environment is a non-production environment
func isDevEnvironment() bool {
	return devEnvironments[os.Getenv("APP_ENV")]
}

// isLocalhost checks if the host is a localhost/loopback/private-network address
func isLocalhost(host string, r *http.Request) bool {
	// Early return if already secured via TLS
	if r.TLS != nil {
		return true
	}

	// Strip port if present, works for IPv4, IPv6 and hostnames
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	// Strip brackets from IPv6 addresses without port e.g. [::1] -> ::1
	host = strings.Trim(host, "[]")

	// Early return for exact matches
	switch host {
	case "localhost", "127.0.0.1", "::1", "0.0.0.0":
		return true
	}

	// Early return for loopback IP range
	if strings.HasPrefix(host, "127.") {
		return true
	}

	// Early return for private IP ranges
	if strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "192.168.") {
		return true
	}

	// Early return for .local domains
	if strings.HasSuffix(host, ".local") {
		return true
	}

	return false
}

// isHTTPS checks if the request is using HTTPS, including proxy-terminated TLS
func isHTTPS(r *http.Request) bool {
	// Early return if HTTPS via TLS
	if r.TLS != nil {
		return true
	}

	// Early return if HTTPS via URL scheme
	if r.URL.Scheme == "https" {
		return true
	}

	// Early return if HTTPS via X-Forwarded-Proto header
	if strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		return true
	}

	// Early return if HTTPS via X-Forwarded-Scheme header
	if strings.EqualFold(r.Header.Get("X-Forwarded-Scheme"), "https") {
		return true
	}

	return false
}

// isRequired reports whether the request should be redirected to HTTPS
func isRequired(r *http.Request) bool {
	// Early return if already HTTPS
	if isHTTPS(r) {
		return false
	}

	// Early return for loopback / private-network hosts
	if isLocalhost(r.Host, r) {
		return false
	}

	// Early return for non-production environments
	if isDevEnvironment() {
		return false
	}

	return true
}

// Config holds optional configuration for the HTTPS redirect middleware.
type Config struct {
	// SkipFunc, if non-nil, is called for every request. When it returns true,
	// the middleware passes the request through without redirecting.
	SkipFunc func(*http.Request) bool
}

// NewHTTPSRedirectMiddleware returns middleware that permanently redirects plain
// HTTP requests to their HTTPS equivalent. Redirection is skipped for:
//   - Requests already using HTTPS (including proxy-terminated TLS)
//   - Loopback / private-network hosts
//   - Non-production environments (APP_ENV ∈ development | local | testing)
func NewHTTPSRedirectMiddleware() rtr.MiddlewareInterface {
	return NewHTTPSRedirectMiddlewareWithConfig(Config{})
}

// NewHTTPSRedirectMiddlewareWithConfig creates middleware with custom configuration.
func NewHTTPSRedirectMiddlewareWithConfig(config Config) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("HTTPS Redirect").
		SetHandler(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !isRequired(r) {
					next.ServeHTTP(w, r)
					return
				}

				if config.SkipFunc != nil && config.SkipFunc(r) {
					next.ServeHTTP(w, r)
					return
				}

				target := "https://" + r.Host + r.URL.RequestURI()
				http.Redirect(w, r, target, http.StatusMovedPermanently)
			})
		})
}
