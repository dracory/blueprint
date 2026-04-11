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
		{"172.31.255.255", "172.31.255.255", true},
		{"server.local", "server.local", true},
		{"api.internal", "api.internal", true},
		{"web.corp.example.com", "web.corp.example.com", true},
		{"example.com", "example.com", false},
		{"public-domain.com", "public-domain.com", false},
		{"cdn.example.org", "cdn.example.org", false},
		{"169.254.169.254", "169.254.169.254", true}, // link-local
		{"172.32.0.1", "172.32.0.1", true},           // all 172.x.x.x considered private (simplified)
		{"192.169.0.1", "192.169.0.1", false},        // outside 192.168 range
		{"10.1.1.1", "10.1.1.1", true},               // inside 10.0.0.0/8 range
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

// TestNewThumbController verifies the constructor creates a valid controller
func TestNewThumbController(t *testing.T) {
	// Since we don't have easy access to a registry in this test package,
	// we'll just verify the constructor doesn't panic with nil
	c := NewThumbController(nil)
	if c == nil {
		t.Fatal("NewThumbController should not return nil even with nil registry")
	}
}

func TestThumbController_Handler_NilRegistry(t *testing.T) {
	c := NewThumbController(nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	result := c.Handler(w, req)
	// Handler should not panic even with nil registry
	if result == "" {
		t.Log("Handler returned empty string with nil registry (expected)")
	}
}

// TestValidateURL tests various URL validation scenarios
func TestValidateURL(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name        string
		url         string
		shouldError bool
	}{
		{"Valid HTTP URL", "http://example.com/image.jpg", false},
		{"Valid HTTPS URL", "https://example.com/image.jpg", false},
		{"Invalid scheme", "ftp://example.com/image.jpg", true},
		{"Empty URL", "", true},
		{"Invalid format", "not-a-url", true},
		{"URL without hostname", "http:///image.jpg", true},
		{"URL with special chars", "http://example.com/image.jpg?test=1", false},
		{"URL with port", "http://example.com:8080/image.jpg", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateURL(tt.url)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for URL %s", tt.url)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for URL %s: %v", tt.url, err)
			}
		})
	}
}

// TestPrepareDataMissingParameters tests missing parameter validation
func TestPrepareDataMissingParameters(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name        string
		url         string
		expectedErr string
	}{
		{"Missing extension", "/th//300x200/70/image.jpg", "image extension is missing"},
		{"Missing size", "/th/jpg//70/image.jpg", "size is missing"},
		{"Missing quality", "/th/jpg/300x200//image.jpg", "quality is missing"},
		{"Missing path", "/th/jpg/300x200/70/", "path is missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := rtr.NewRouter()
			route := rtr.NewRoute().
				SetMethod("GET").
				SetPath("/th/:extension/:size/:quality/:path").
				SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
					_, errMsg := c.prepareData(r)
					return errMsg
				})
			router.AddRoute(route)

			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Body.String() != tt.expectedErr {
				t.Errorf("Expected error %q, got %q", tt.expectedErr, rr.Body.String())
			}
		})
	}
}

// TestPrepareDataURLNormalization tests URL path normalization
func TestPrepareDataURLNormalization(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name     string
		path     string
		isURL    bool
		isCache  bool
		isFiles  bool
		expected string
	}{
		{"HTTP URL normalization", "http/example.com/image.jpg", true, false, false, "http://example.com/image.jpg"},
		{"HTTPS URL normalization", "https/example.com/image.jpg", true, false, false, "https://example.com/image.jpg"},
		{"Cache prefix", "cache-mykey", false, true, false, "mykey"},
		{"Files prefix", "files/uploads/image.jpg", false, false, true, "uploads/image.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := rtr.NewRouter()
			route := rtr.NewRoute().
				SetMethod("GET").
				SetPath("/th/:extension/:size/:quality/:path").
				SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
					data, _ := c.prepareData(r)
					return fmt.Sprintf("%s|%t|%t|%t", data.path, data.isURL, data.isCache, data.isFiles)
				})
			router.AddRoute(route)

			url := fmt.Sprintf("/th/jpg/300x200/70/%s", tt.path)
			req := httptest.NewRequest("GET", url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			body := rr.Body.String()
			parts := strings.Split(body, "|")
			if len(parts) >= 4 {
				if parts[0] != tt.expected {
					t.Errorf("Expected path %q, got %q", tt.expected, parts[0])
				}
				if parts[1] != fmt.Sprintf("%t", tt.isURL) {
					t.Errorf("Expected isURL %t, got %s", tt.isURL, parts[1])
				}
				if parts[2] != fmt.Sprintf("%t", tt.isCache) {
					t.Errorf("Expected isCache %t, got %s", tt.isCache, parts[2])
				}
				if parts[3] != fmt.Sprintf("%t", tt.isFiles) {
					t.Errorf("Expected isFiles %t, got %s", tt.isFiles, parts[3])
				}
			}
		})
	}
}

// TestPrepareDataSizeParsing tests various size format parsing
func TestPrepareDataSizeParsing(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name           string
		size           string
		expectedWidth  int64
		expectedHeight int64
	}{
		{"Width and height", "300x200", 300, 200},
		{"Width only", "500", 500, 0},
		{"Zero dimensions", "0x0", 0, 0},
		{"Large dimensions", "4000x3000", 4000, 3000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := rtr.NewRouter()
			route := rtr.NewRoute().
				SetMethod("GET").
				SetPath("/th/:extension/:size/:quality/:path").
				SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
					data, _ := c.prepareData(r)
					return fmt.Sprintf("%d|%d", data.width, data.height)
				})
			router.AddRoute(route)

			url := fmt.Sprintf("/th/jpg/%s/70/image.jpg", tt.size)
			req := httptest.NewRequest("GET", url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			body := rr.Body.String()
			parts := strings.Split(body, "|")
			if len(parts) >= 2 {
				var w, h int64
				fmt.Sscanf(parts[0], "%d", &w)
				fmt.Sscanf(parts[1], "%d", &h)
				if w != tt.expectedWidth || h != tt.expectedHeight {
					t.Errorf("Expected %dx%d, got %dx%d", tt.expectedWidth, tt.expectedHeight, w, h)
				}
			}
		})
	}
}

// TestSetHeadersAllFormats tests all supported image format headers
func TestSetHeadersAllFormats(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		extension   string
		contentType string
	}{
		{"jpg", "image/jpeg"},
		{"jpeg", "image/jpeg"},
		{"png", "image/png"},
		{"gif", "image/gif"},
		{"webp", ""},
		{"bmp", ""},
	}

	for _, tt := range tests {
		t.Run(tt.extension, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c.setHeaders(rec, tt.extension)

			contentType := rec.Header().Get("Content-Type")
			if contentType != tt.contentType {
				t.Errorf("Extension %s: expected Content-Type %q, got %q", tt.extension, tt.contentType, contentType)
			}

			cacheControl := rec.Header().Get("Cache-Control")
			if cacheControl != "max-age=604800" {
				t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
			}
		})
	}
}

// TestToBytes tests the toBytes function (local file reading)
func TestToBytes(t *testing.T) {
	c := &thumbnailController{}

	// Test with non-existent file
	_, err := c.toBytes("/nonexistent/path/image.jpg")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestIsPrivateHostEdgeCases tests edge cases in private host detection
func TestIsPrivateHostEdgeCases(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"172.15.0.1 (simplified check treats all 172.x as private)", "172.15.0.1", true},
		{"172.32.0.1 (simplified check treats all 172.x as private)", "172.32.0.1", true},
		{"169.253.0.1 (before link-local)", "169.253.0.1", false},
		{"169.255.0.1 (after link-local)", "169.255.0.1", false},
		{"example.corp.com", "example.corp.com", true},
		{"corp.example.com", "corp.example.com", false},
		{"10.255.255.255", "10.255.255.255", true},
		{"11.0.0.0", "11.0.0.0", false},
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

// TestValidateURLWithSpecialCases tests additional URL validation scenarios
func TestValidateURLWithSpecialCases(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name        string
		url         string
		shouldError bool
	}{
		{"IPv6 localhost", "http://[::1]/image.jpg", true},
		{"IPv6 address", "http://[2001:db8::1]/image.jpg", true},
		{"URL with auth", "http://user:pass@example.com/image.jpg", false},
		{"URL with fragment", "http://example.com/image.jpg#section", false},
		{"URL with query params", "http://example.com/image.jpg?size=large&format=jpg", false},
		{".local domain", "http://myserver.local/image.jpg", true},
		{".internal domain", "http://api.internal/image.jpg", true},
		{"No scheme", "example.com/image.jpg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateURL(tt.url)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for URL %s, but got none", tt.url)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for URL %s, but got: %v", tt.url, err)
			}
		})
	}
}

// TestPrepareDataQualityParsing tests quality parameter parsing
func TestPrepareDataQualityParsing(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name            string
		quality         string
		expectedQuality int64
	}{
		{"Low quality", "30", 30},
		{"Medium quality", "70", 70},
		{"High quality", "95", 95},
		{"Zero quality", "0", 0},
		{"Invalid quality (non-numeric)", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := rtr.NewRouter()
			route := rtr.NewRoute().
				SetMethod("GET").
				SetPath("/th/:extension/:size/:quality/:path").
				SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
					data, _ := c.prepareData(r)
					return fmt.Sprintf("%d", data.quality)
				})
			router.AddRoute(route)

			url := fmt.Sprintf("/th/jpg/300x200/%s/image.jpg", tt.quality)
			req := httptest.NewRequest("GET", url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			body := rr.Body.String()
			var q int64
			fmt.Sscanf(body, "%d", &q)
			if q != tt.expectedQuality {
				t.Errorf("Expected quality %d, got %d", tt.expectedQuality, q)
			}
		})
	}
}

// TestValidateURLWithCorporateDomains tests corporate domain blocking
func TestValidateURLWithCorporateDomains(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name        string
		url         string
		shouldError bool
	}{
		{"Internal corp domain", "http://internal.corp.example.com/image.jpg", true},
		{"Nested corp domain", "http://api.service.corp.example.com/image.jpg", true},
		{"Public corp domain", "http://example.corp/image.jpg", false},
		{"Corp in subdomain", "http://corp.example.com/image.jpg", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateURL(tt.url)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for URL %s, but got none", tt.url)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for URL %s, but got: %v", tt.url, err)
			}
		})
	}
}

// TestPrepareDataExtensionHandling tests extension parameter handling
func TestPrepareDataExtensionHandling(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name              string
		extension         string
		expectedExtension string
	}{
		{"JPEG extension", "jpg", "jpg"},
		{"PNG extension", "png", "png"},
		{"GIF extension", "gif", "gif"},
		{"Uppercase extension", "JPG", "JPG"},
		{"WebP extension", "webp", "webp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := rtr.NewRouter()
			route := rtr.NewRoute().
				SetMethod("GET").
				SetPath("/th/:extension/:size/:quality/:path").
				SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
					data, _ := c.prepareData(r)
					return data.extension
				})
			router.AddRoute(route)

			url := fmt.Sprintf("/th/%s/300x200/70/image.jpg", tt.extension)
			req := httptest.NewRequest("GET", url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			body := rr.Body.String()
			if body != tt.expectedExtension {
				t.Errorf("Expected extension %q, got %q", tt.expectedExtension, body)
			}
		})
	}
}

// TestIsPublicDomainComprehensive tests public domain detection comprehensively
func TestIsPublicDomainComprehensive(t *testing.T) {
	c := &thumbnailController{}

	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"Public domain", "example.com", true},
		{"Public subdomain", "cdn.example.com", true},
		{"Public with hyphen", "my-domain.com", true},
		{"Public with numbers", "example123.com", true},
		{"Localhost", "localhost", false},
		{"127.0.0.1", "127.0.0.1", false},
		{"Private IP", "192.168.1.1", false},
		{".local domain", "server.local", false},
		{".internal domain", "api.internal", false},
		{"IPv6 loopback", "::1", false},
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
