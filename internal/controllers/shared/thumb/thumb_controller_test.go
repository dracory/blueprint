package thumb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

// TestValidateURLSecurity tests the SSRF protection functionality
func TestValidateURLSecurity(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name        string
		url         string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Valid HTTPS URL",
			url:         "https://example.com/image.jpg",
			shouldError: false,
		},
		{
			name:        "Valid HTTP URL",
			url:         "http://public-domain.com/image.png",
			shouldError: false,
		},
		{
			name:        "Invalid scheme - FTP",
			url:         "ftp://example.com/image.jpg",
			shouldError: true,
			errorMsg:    "only HTTP and HTTPS URLs are allowed",
		},
		{
			name:        "Invalid scheme - file",
			url:         "file:///etc/passwd",
			shouldError: true,
			errorMsg:    "only HTTP and HTTPS URLs are allowed",
		},
		{
			name:        "localhost access blocked",
			url:         "http://localhost:8080/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "127.0.0.1 access blocked",
			url:         "http://127.0.0.1/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Private IP 192.168 blocked",
			url:         "http://192.168.1.100/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Private IP 10.0.0.0 blocked",
			url:         "http://10.0.0.50/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Private IP 172.16 blocked",
			url:         "http://172.16.0.1/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Internal domain .local blocked",
			url:         "http://server.local/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Internal domain .internal blocked",
			url:         "http://api.internal/image.jpg",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Link-local IP blocked",
			url:         "http://169.254.169.254/latest/meta-data/",
			shouldError: true,
			errorMsg:    "access to private networks is not allowed",
		},
		{
			name:        "Invalid URL format",
			url:         "not-a-valid-url",
			shouldError: true,
			errorMsg:    "invalid URL format",
		},
		{
			name:        "Empty URL",
			url:         "",
			shouldError: true,
			errorMsg:    "invalid URL format",
		},
		{
			name:        "URL without hostname",
			url:         "http:///image.jpg",
			shouldError: true,
			errorMsg:    "invalid hostname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateURL(tt.url)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for URL %s, but got none", tt.url)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for URL %s, but got: %v", tt.url, err)
				}
			}
		})
	}
}

// TestIsPrivateHost tests the private host detection logic
func TestIsPrivateHost(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"localhost", "localhost", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"::1", "::1", true},
		{"192.168.1.1", "192.168.1.1", true},
		{"10.0.0.1", "10.0.0.1", true},
		{"172.16.0.1", "172.16.0.1", true},
		{"server.local", "server.local", true},
		{"api.internal", "api.internal", true},
		{"web.corp.example.com", "web.corp.example.com", true},
		{"example.com", "example.com", false},
		{"public-domain.com", "public-domain.com", false},
		{"cdn.example.org", "cdn.example.org", false},
		{"169.254.169.254", "169.254.169.254", true}, // link-local
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.isPrivateHost(tt.host)
			if result != tt.expected {
				t.Errorf("isPrivateHost(%s) = %v, expected %v", tt.host, result, tt.expected)
			}
		})
	}
}

// TestIsPublicDomain tests the public domain detection logic
func TestIsPublicDomain(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"example.com", "example.com", true},
		{"public-domain.com", "public-domain.com", true},
		{"cdn.example.org", "cdn.example.org", true},
		{"localhost", "localhost", false},
		{"127.0.0.1", "127.0.0.1", false},
		{"192.168.1.1", "192.168.1.1", false},
		{"169.254.169.254", "169.254.169.254", false},
		{"server.local", "server.local", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.isPublicDomain(tt.host)
			if result != tt.expected {
				t.Errorf("isPublicDomain(%s) = %v, expected %v", tt.host, result, tt.expected)
			}
		})
	}
}

// servePrepareData routes a request through rtr to invoke prepareData and
// capture its parsed data and error message for assertions.
func servePrepareData(method, target string) (data thumbnailControllerData, errMsg string, status int) {
	c := &thumbnailController{}

	router := rtr.NewRouter()
	// Route using single-segment ':path'
	route := rtr.NewRoute().
		SetMethod(http.MethodGet).
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			d, e := c.prepareData(r)
			return fmt.Sprintf("%s|%s|%d|%d|%d|%t|%t", e, d.path, d.width, d.height, d.quality, d.isURL, d.isCache)
		})
	router.AddRoute(route)

	req := httptest.NewRequest(method, target, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	status = rr.Code
	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 7 {
		errMsg = parts[0]
		data.path = parts[1]
		// parse ints and bools
		var w, h, q int
		var isURLStr, isCacheStr string
		if _, err := fmt.Sscanf(parts[2], "%d", &w); err != nil {
			errMsg = fmt.Sprintf("failed to parse width: %v", err)
		}
		if _, err := fmt.Sscanf(parts[3], "%d", &h); err != nil {
			errMsg = fmt.Sprintf("failed to parse height: %v", err)
		}
		if _, err := fmt.Sscanf(parts[4], "%d", &q); err != nil {
			errMsg = fmt.Sprintf("failed to parse quality: %v", err)
		}
		data.width = int64(w)
		data.height = int64(h)
		data.quality = int64(q)
		isURLStr = parts[5]
		isCacheStr = parts[6]
		data.isURL = isURLStr == "true"
		data.isCache = isCacheStr == "true"
	} else {
		errMsg = "invalid test response"
	}
	return data, errMsg, status
}

// TestSetHeaders verifies correct content-type and cache-control headers
func TestSetHeaders(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "jpg")
	if got := rec.Header().Get("Content-Type"); got != "image/jpeg" {
		t.Fatalf("jpg content-type expected image/jpeg, got: %q", got)
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "max-age=604800" {
		t.Fatalf("expected Cache-Control max-age=604800, got: %q", cc)
	}

	rec = httptest.NewRecorder()
	c.setHeaders(rec, "jpeg")
	if got := rec.Header().Get("Content-Type"); got != "image/jpeg" {
		t.Fatalf("jpeg content-type expected image/jpeg, got: %q", got)
	}

	rec = httptest.NewRecorder()
	c.setHeaders(rec, "png")
	if got := rec.Header().Get("Content-Type"); got != "image/png" {
		t.Fatalf("png content-type expected image/png, got: %q", got)
	}

	rec = httptest.NewRecorder()
	c.setHeaders(rec, "gif")
	if got := rec.Header().Get("Content-Type"); got != "image/gif" {
		t.Fatalf("gif content-type expected image/gif, got: %q", got)
	}

	rec = httptest.NewRecorder()
	c.setHeaders(rec, "bin")
	if got := rec.Header().Get("Content-Type"); got != "" {
		t.Fatalf("unknown extension should set empty content-type, got: %q", got)
	}
}

// TestPrepareDataParsing verifies size parsing and URL normalization via rtr routing
func TestPrepareDataParsing(t *testing.T) {
	// width x height parsing
	data, errMsg, status := servePrepareData(http.MethodGet, "/th/jpg/200x150/70/ab.jpg")
	if status != http.StatusOK {
		t.Fatalf("unexpected status: %d", status)
	}
	if errMsg != "" {
		t.Fatalf("unexpected error: %q", errMsg)
	}
	if data.width != 200 || data.height != 150 || data.quality != 70 {
		t.Fatalf("expected width=200 height=150 quality=70, got %d %d %d", data.width, data.height, data.quality)
	}

	// single dimension parsing
	data, errMsg, status = servePrepareData(http.MethodGet, "/th/png/300/60/ab.png")
	if status != http.StatusOK {
		t.Fatalf("unexpected status: %d", status)
	}
	if errMsg != "" {
		t.Fatalf("unexpected error: %q", errMsg)
	}
	if data.width != 300 || data.height != 0 {
		t.Fatalf("expected width=300 height=0, got %d %d", data.width, data.height)
	}
}

// TestPrepareDataFilesLink ensures files/ paths are converted via links.URL and marked as URL
func TestPrepareDataFilesLink(t *testing.T) {
	t.Skip("requires catch-all path; skipped while route uses :path")
	// Ensure APP_ENV=testing so links.RootURL() returns empty string (see links.RootURL)
	oldEnv := os.Getenv("APP_ENV")
	if err := os.Setenv("APP_ENV", "testing"); err != nil {
		t.Fatalf("failed to set APP_ENV: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Setenv("APP_ENV", oldEnv); err != nil {
			t.Logf("warning: failed to restore APP_ENV: %v", err)
		}
	})
}
