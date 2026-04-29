package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	// Create test config with development environment
	testConfig := config.New()
	testConfig.SetAppEnv(config.APP_ENVIRONMENT_DEVELOPMENT)
	testConfig.SetDatabaseDriver("sqlite") // Required by registry builder
	testConfig.SetDatabaseName(":memory:") // Required by registry builder
	registry := testutils.Setup(testutils.WithCfg(testConfig))

	cspParts := []string{
		"default-src 'self';",
		"script-src 'unsafe-inline' 'unsafe-hashes' 'unsafe-eval' 'self'",
		"https://cdn.jsdelivr.net",
		"http://cdn.jsdelivr.net",
		"https://unpkg.com",
		"https://www.statcounter.com",
		"https://code.jquery.com",
		"https://cdn.datatables.net",
		"https://cdnjs.cloudflare.com",
		"http://cdnjs.cloudflare.com",
		"https://www.googletagmanager.com",
		"https://www.statcounter.com;",
		"style-src 'unsafe-inline' 'unsafe-hashes' 'self'",
		"https://cdn.jsdelivr.net",
		"https://maxcdn.bootstrapcdn.com",
		"https://cdnjs.cloudflare.com",
		"http://cdnjs.cloudflare.com",
		"https://fonts.googleapis.com",
		"https://unpkg.com",
		"https://code.jquery.com",
		"https://cdn.datatables.net",
		"https://cdnjs.cloudflare.com;",
		"font-src 'self'",
		"https://cdn.jsdelivr.net",
		"https://fonts.googleapis.com",
		"https://fonts.gstatic.com",
		"https://cdnjs.cloudflare.com",
		"http://cdnjs.cloudflare.com",
		"https://maxcdn.bootstrapcdn.com;",
		"img-src 'self' data:",
		"https://sfs.ams3.digitaloceanspaces.com",
		"https://lesichkov.ams3.digitaloceanspaces.com",
		"https://provedexpert.gitlab.io;",
		"connect-src 'self'",
		"https://cdnjs.cloudflare.com",
		"http://cdnjs.cloudflare.com",
	}

	tests := []struct {
		header        string
		expectedValue string
	}{
		// HSTS is disabled in development
		{"X-Frame-Options", "DENY"},
		{"X-Content-Type-Options", "nosniff"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{
			"Content-Security-Policy",
			strings.Join(cspParts, " "),
		},
	}

	req := httptest.NewRequest("GET", "https://example.com", nil)
	rr := httptest.NewRecorder()
	handler := NewSecurityHeadersMiddleware(registry).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			if value := rr.Header().Get(tt.header); value != tt.expectedValue {
				t.Errorf("header %s: got %q want %q", tt.header, value, tt.expectedValue)
			}
		})
	}
}
