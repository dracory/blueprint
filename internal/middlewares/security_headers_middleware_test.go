package middlewares

import (
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"strings"
	"testing"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	// Create test config with development environment
	testConfig := config.New()
	testConfig.SetAppEnv(config.APP_ENVIRONMENT_DEVELOPMENT)
	testConfig.SetDatabaseDriver("sqlite") // Required by app builder
	testConfig.SetDatabaseName(":memory:") // Required by app builder
	app := testutils.Setup(testutils.WithCfg(testConfig))

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
		"https://www.statcounter.com",
		"https://cdn.tiny.cloud;",
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
		"https://www.statcounter.com",
	}

	req := httptest.NewRequest("GET", "https://example.com", nil)
	rr := httptest.NewRecorder()
	handler := NewSecurityHeadersMiddleware(app).GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if value := rr.Header().Get("X-Frame-Options"); value != "DENY" {
		t.Errorf("header X-Frame-Options: got %q want DENY", value)
	}
	if value := rr.Header().Get("X-Content-Type-Options"); value != "nosniff" {
		t.Errorf("header X-Content-Type-Options: got %q want nosniff", value)
	}
	if value := rr.Header().Get("Referrer-Policy"); value != "strict-origin-when-cross-origin" {
		t.Errorf("header Referrer-Policy: got %q want strict-origin-when-cross-origin", value)
	}
	expectedCSP := strings.Join(cspParts, " ")
	if value := rr.Header().Get("Content-Security-Policy"); value != expectedCSP {
		t.Errorf("header Content-Security-Policy: got %q want %q", value, expectedCSP)
	}
}
