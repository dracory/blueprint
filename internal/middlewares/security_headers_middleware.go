package middlewares

import (
	"os"
	"strings"

	"github.com/dracory/rtr"
	"github.com/dracory/rtr/middlewares"
)

// NewSecurityHeadersMiddleware creates middleware that sets security headers
// Uses RTR security middleware with project-specific configuration
func NewSecurityHeadersMiddleware() rtr.MiddlewareInterface {
	isDevelopment := os.Getenv("APP_ENV") == "development"

	config := &middlewares.SecurityHeadersConfig{
		CSP: &middlewares.CSPConfig{
			Enabled:    true,
			DefaultSrc: []string{"'self'"},
			ScriptSrc:  getScriptSources(isDevelopment),
			StyleSrc:   getStyleSources(isDevelopment),
			FontSrc: []string{
				"'self'",
				"https://cdn.jsdelivr.net",
				"https://fonts.googleapis.com",
				"https://fonts.gstatic.com",
				"https://cdnjs.cloudflare.com",
				"https://maxcdn.bootstrapcdn.com",
			},
			ImgSrc: []string{
				"'self'",
				"data:",
				"https://sfs.ams3.digitaloceanspaces.com",
				"https://lesichkov.ams3.digitaloceanspaces.com",
			},
			UpgradeInsecureRequests: !isDevelopment,
		},
		HSTS: &middlewares.HSTSConfig{
			Enabled:           !isDevelopment,
			MaxAge:            31536000,
			IncludeSubDomains: !isDevelopment,
			Preload:           !isDevelopment,
		},
		FrameOptions: &middlewares.FrameOptionsConfig{
			Enabled: true,
			Option:  "DENY",
		},
		ContentTypeNosniff: true,
		XSSProtection: &middlewares.XSSProtectionConfig{
			Enabled: true,
			Mode:    "block",
		},
		ReferrerPolicy: "strict-origin-when-cross-origin",
		PermissionsPolicy: map[string][]string{
			"geolocation": {},
			"microphone":  {},
			"camera":      {},
		},
		CustomHeaders: getCustomHeaders(isDevelopment),
	}

	return middlewares.NewSecurityHeadersMiddleware(config)
}

// NewCustomSecurityHeadersMiddleware allows full customization
func NewCustomSecurityHeadersMiddleware(config *middlewares.SecurityHeadersConfig) rtr.MiddlewareInterface {
	return middlewares.NewSecurityHeadersMiddleware(config)
}

// getScriptSources returns script sources based on environment
func getScriptSources(isDevelopment bool) []string {
	sources := []string{
		"'self'",
		"https://cdn.jsdelivr.net",
		"https://unpkg.com",
		"https://code.jquery.com",
		"https://cdnjs.cloudflare.com",
		"https://www.googletagmanager.com",
		"https://www.statcounter.com",
	}

	if isDevelopment {
		// Allow unsafe-inline and eval in development
		sources = append([]string{"'unsafe-inline'", "'unsafe-eval'"}, sources...)
	}

	return sources
}

// getStyleSources returns style sources based on environment
func getStyleSources(isDevelopment bool) []string {
	sources := []string{
		"'self'",
		"https://cdn.jsdelivr.net",
		"https://maxcdn.bootstrapcdn.com",
		"https://cdnjs.cloudflare.com",
		"https://fonts.googleapis.com",
		"https://unpkg.com",
	}

	if isDevelopment {
		// Allow unsafe-inline in development
		sources = append([]string{"'unsafe-inline'"}, sources...)
	}

	return sources
}

// getCustomHeaders returns custom headers based on environment
func getCustomHeaders(isDevelopment bool) map[string]string {
	headers := make(map[string]string)

	if !isDevelopment {
		// Production-only headers
		headers["X-Content-Type-Options"] = "nosniff"
		headers["X-Frame-Options"] = "DENY"
	}

	// Add project-specific custom headers from environment
	if customHeaders := os.Getenv("CUSTOM_SECURITY_HEADERS"); customHeaders != "" {
		for _, header := range strings.Split(customHeaders, ",") {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return headers
}
