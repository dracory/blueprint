package middlewares

import (
	"net/http"
	"os"

	"github.com/dracory/rtr"
	"github.com/dracory/rtr/middlewares"
)

// NewHTTPSRedirectMiddleware creates middleware that redirects HTTP requests to HTTPS
// Uses RTR security middleware with project-specific configuration
func NewHTTPSRedirectMiddleware() rtr.MiddlewareInterface {
	config := &middlewares.HTTPSRedirectConfig{
		SkipLocalhost: os.Getenv("APP_ENV") == "development",
		TrustedProxies: []string{
			"127.0.0.1",
			"::1",
		},
		CustomSkipFunc: func(r *http.Request) bool {
			// Project-specific skip logic
			// Skip HTTPS redirect for health checks in development
			if os.Getenv("APP_ENV") == "development" && r.URL.Path == "/health" {
				return true
			}

			// Skip for API endpoints that might be called internally
			if r.URL.Path == "/api/internal/webhook" {
				return true
			}

			return false
		},
	}

	return middlewares.NewHTTPSRedirectMiddleware(config)
}

// NewCustomHTTPSRedirectMiddleware allows full customization
func NewCustomHTTPSRedirectMiddleware(config *middlewares.HTTPSRedirectConfig) rtr.MiddlewareInterface {
	return middlewares.NewHTTPSRedirectMiddleware(config)
}
