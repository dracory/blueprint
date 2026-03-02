package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPSRedirectMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		wantStatus     int
		wantLocation   string
		isLocalRequest bool
	}{
		{"redirects HTTP to HTTPS", "http://example.com/foo", http.StatusMovedPermanently, "https://example.com/foo", false},
		{"no redirect for HTTPS", "https://example.com/bar", http.StatusOK, "", false},
		{"no redirect for localhost", "http://localhost:8080", http.StatusOK, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			if tt.isLocalRequest {
				req.Host = "localhost"
			}

			rr := httptest.NewRecorder()
			handler := NewHTTPSRedirectMiddleware().GetHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			if tt.wantLocation != "" {
				if location := rr.Header().Get("Location"); location != tt.wantLocation {
					t.Errorf("handler returned wrong location header: got %v want %v", location, tt.wantLocation)
				}
			}
		})
	}
}
