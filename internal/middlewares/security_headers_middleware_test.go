package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	tests := []struct {
		header        string
		expectedValue string
	}{
		{"Strict-Transport-Security", "max-age=31536000; includeSubDomains"},
		{"X-Frame-Options", "DENY"},
		{"X-Content-Type-Options", "nosniff"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{"Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://unpkg.com https://code.jquery.com https://cdnjs.cloudflare.com https://www.googletagmanager.com https://www.statcounter.com; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://maxcdn.bootstrapcdn.com https://cdnjs.cloudflare.com https://fonts.googleapis.com; font-src 'self' https://cdn.jsdelivr.net https://fonts.googleapis.com https://fonts.gstatic.com https://cdnjs.cloudflare.com https://maxcdn.bootstrapcdn.com; img-src 'self' data: https://sfs.ams3.digitaloceanspaces.com https://lesichkov.ams3.digitaloceanspaces.com"},
	}

	req := httptest.NewRequest("GET", "https://example.com", nil)
	rr := httptest.NewRecorder()
	handler := NewSecurityHeadersMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
