package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func writeMaintenanceFile(t *testing.T, path string, state MaintenanceState) {
	t.Helper()
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("failed to marshal state: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write maintenance file: %v", err)
	}
}

func TestMaintenanceMiddleware_NoFile_PassesThrough(t *testing.T) {
	mw := &maintenanceMiddleware{
		filePath: "test_maintenance_nonexistent.json",
		cacheDur: 50 * time.Millisecond,
	}

	called := false
	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called when no maintenance file exists")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestMaintenanceMiddleware_FileExists_Returns503(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		Message:           "Down for maintenance",
		RetryAfterSeconds: 60,
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	called := false
	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if called {
		t.Fatal("expected next handler NOT to be called when maintenance is active")
	}
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") != "60" {
		t.Fatalf("expected Retry-After header '60', got '%s'", w.Header().Get("Retry-After"))
	}
}

func TestMaintenanceMiddleware_ExcludedPath_PassesThrough(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		ExcludePaths: []string{"/admin/*"},
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	called := false
	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called for excluded path")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for excluded path, got %d", w.Code)
	}
}

func TestMaintenanceMiddleware_ExcludedIP_PassesThrough(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		ExcludeIPs: []string{"203.0.113.5"},
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	called := false
	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "203.0.113.5:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called for excluded IP")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for excluded IP, got %d", w.Code)
	}
}

func TestMaintenanceMiddleware_ExcludedIP_XForwardedFor(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		ExcludeIPs: []string{"198.51.100.10"},
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	called := false
	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "198.51.100.10, 10.0.0.1")
	req.RemoteAddr = "10.0.0.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called for excluded IP via X-Forwarded-For")
	}
}

func TestMaintenanceMiddleware_NonExcludedIP_Returns503(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		ExcludeIPs: []string{"203.0.113.5"},
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	called := false
	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.99:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if called {
		t.Fatal("expected next handler NOT to be called for non-excluded IP")
	}
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", w.Code)
	}
}

func TestMaintenanceMiddleware_RetryAfterHeader(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		RetryAfterSeconds: 120,
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Retry-After") != "120" {
		t.Fatalf("expected Retry-After '120', got '%s'", w.Header().Get("Retry-After"))
	}
}

func TestMaintenanceMiddleware_CustomMessage(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		Message: "Database migration in progress",
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 50 * time.Millisecond,
	}

	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "Database migration in progress") {
		t.Fatalf("expected body to contain custom message, got: %s", body)
	}
}

func TestMaintenanceMiddleware_FileCache(t *testing.T) {
	path := "test_maintenance_state.json"
	defer os.Remove(path)

	writeMaintenanceFile(t, path, MaintenanceState{
		Message: "First message",
	})

	mw := &maintenanceMiddleware{
		filePath: path,
		cacheDur: 5 * time.Second,
	}

	handler := mw.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1:12345"

	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req)
	body1 := w1.Body.String()

	if !strings.Contains(body1, "First message") {
		t.Fatalf("expected first message, got: %s", body1)
	}

	writeMaintenanceFile(t, path, MaintenanceState{
		Message: "Second message",
	})

	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)
	body2 := w2.Body.String()

	if !strings.Contains(body2, "First message") {
		t.Fatalf("expected cached first message (within cache window), got: %s", body2)
	}
}

func TestIsIPExcluded(t *testing.T) {
	tests := []struct {
		ip         string
		excludeIPs []string
		expected   bool
	}{
		{"203.0.113.5", []string{"203.0.113.5"}, true},
		{"203.0.113.5", []string{"203.0.113.6"}, false},
		{"203.0.113.5", []string{"203.0.113.5", "198.51.100.10"}, true},
		{"203.0.113.5", []string{}, false},
		{"203.0.113.5", nil, false},
		{" 203.0.113.5 ", []string{"203.0.113.5"}, false},
	}

	for _, tt := range tests {
		result := isIPExcluded(tt.ip, tt.excludeIPs)
		if result != tt.expected {
			t.Errorf("isIPExcluded(%q, %v) = %v, expected %v", tt.ip, tt.excludeIPs, result, tt.expected)
		}
	}
}

func TestIsPathExcluded(t *testing.T) {
	tests := []struct {
		path     string
		patterns []string
		expected bool
	}{
		{"/admin/dashboard", []string{"/admin/*"}, true},
		{"/admin", []string{"/admin/*"}, true},
		{"/api/health", []string{"/api/health"}, true},
		{"/api/health", []string{"/api/*"}, true},
		{"/blog", []string{"/admin/*"}, false},
		{"/blog/post", []string{"/blog*"}, true},
		{"/", []string{"/admin/*"}, false},
		{"/admin/users", []string{"/admin"}, false},
	}

	for _, tt := range tests {
		result := isPathExcluded(tt.path, tt.patterns)
		if result != tt.expected {
			t.Errorf("isPathExcluded(%q, %v) = %v, expected %v", tt.path, tt.patterns, result, tt.expected)
		}
	}
}
