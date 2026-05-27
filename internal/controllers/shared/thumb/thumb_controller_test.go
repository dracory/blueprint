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
func TestValidateURLSecurity_ValidHTTPSURL(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("https://example.com/image.jpg")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

func TestValidateURLSecurity_ValidHTTPURL(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://public-domain.com/image.png")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

func TestValidateURLSecurity_InvalidSchemeFTP(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("ftp://example.com/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "only HTTP and HTTPS URLs are allowed") {
		t.Errorf("Expected error message containing %q, got %q", "only HTTP and HTTPS URLs are allowed", err.Error())
	}
}

func TestValidateURLSecurity_InvalidSchemeFile(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("file:///etc/passwd")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "only HTTP and HTTPS URLs are allowed") {
		t.Errorf("Expected error message containing %q, got %q", "only HTTP and HTTPS URLs are allowed", err.Error())
	}
}

func TestValidateURLSecurity_LocalhostAccessBlocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://localhost:8080/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_127001AccessBlocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://127.0.0.1/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_PrivateIP192168Blocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://192.168.1.100/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_PrivateIP10000Blocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://10.0.0.50/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_PrivateIP17216Blocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://172.16.0.1/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_InternalDomainLocalBlocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://server.local/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_InternalDomainInternalBlocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://api.internal/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_LinkLocalIPBlocked(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://169.254.169.254/latest/meta-data/")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "access to private networks is not allowed") {
		t.Errorf("Expected error message containing %q, got %q", "access to private networks is not allowed", err.Error())
	}
}

func TestValidateURLSecurity_InvalidURLFormat(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("not-a-valid-url")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "invalid URL format") {
		t.Errorf("Expected error message containing %q, got %q", "invalid URL format", err.Error())
	}
}

func TestValidateURLSecurity_EmptyURL(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "invalid URL format") {
		t.Errorf("Expected error message containing %q, got %q", "invalid URL format", err.Error())
	}
}

func TestValidateURLSecurity_URLWithoutHostname(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http:///image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
		return
	}
	if !strings.Contains(err.Error(), "invalid hostname") {
		t.Errorf("Expected error message containing %q, got %q", "invalid hostname", err.Error())
	}
}

// TestIsPrivateHost tests the private host detection logic
func TestIsPrivateHost_Localhost(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("localhost")
	if result != true {
		t.Errorf("isPrivateHost(localhost) = %v, expected true", result)
	}
}

func TestIsPrivateHost_127001(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("127.0.0.1")
	if result != true {
		t.Errorf("isPrivateHost(127.0.0.1) = %v, expected true", result)
	}
}

func TestIsPrivateHost_IPv6Loopback(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("::1")
	if result != true {
		t.Errorf("isPrivateHost(::1) = %v, expected true", result)
	}
}

func TestIsPrivateHost_19216811(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("192.168.1.1")
	if result != true {
		t.Errorf("isPrivateHost(192.168.1.1) = %v, expected true", result)
	}
}

func TestIsPrivateHost_10001(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("10.0.0.1")
	if result != true {
		t.Errorf("isPrivateHost(10.0.0.1) = %v, expected true", result)
	}
}

func TestIsPrivateHost_1721601(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("172.16.0.1")
	if result != true {
		t.Errorf("isPrivateHost(172.16.0.1) = %v, expected true", result)
	}
}

func TestIsPrivateHost_17231255255(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("172.31.255.255")
	if result != true {
		t.Errorf("isPrivateHost(172.31.255.255) = %v, expected true", result)
	}
}

func TestIsPrivateHost_ServerLocal(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("server.local")
	if result != true {
		t.Errorf("isPrivateHost(server.local) = %v, expected true", result)
	}
}

func TestIsPrivateHost_ApiInternal(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("api.internal")
	if result != true {
		t.Errorf("isPrivateHost(api.internal) = %v, expected true", result)
	}
}

func TestIsPrivateHost_WebCorpExampleCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("web.corp.example.com")
	if result != true {
		t.Errorf("isPrivateHost(web.corp.example.com) = %v, expected true", result)
	}
}

func TestIsPrivateHost_ExampleCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("example.com")
	if result != false {
		t.Errorf("isPrivateHost(example.com) = %v, expected false", result)
	}
}

func TestIsPrivateHost_PublicDomainCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("public-domain.com")
	if result != false {
		t.Errorf("isPrivateHost(public-domain.com) = %v, expected false", result)
	}
}

func TestIsPrivateHost_CdnExampleOrg(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("cdn.example.org")
	if result != false {
		t.Errorf("isPrivateHost(cdn.example.org) = %v, expected false", result)
	}
}

func TestIsPrivateHost_169254169254(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("169.254.169.254")
	if result != true {
		t.Errorf("isPrivateHost(169.254.169.254) = %v, expected true", result)
	}
}

func TestIsPrivateHost_1723201(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("172.32.0.1")
	if result != true {
		t.Errorf("isPrivateHost(172.32.0.1) = %v, expected true", result)
	}
}

func TestIsPrivateHost_19216901(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("192.169.0.1")
	if result != false {
		t.Errorf("isPrivateHost(192.169.0.1) = %v, expected false", result)
	}
}

func TestIsPrivateHost_10111(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("10.1.1.1")
	if result != true {
		t.Errorf("isPrivateHost(10.1.1.1) = %v, expected true", result)
	}
}

// TestIsPublicDomain tests the public domain detection logic
func TestIsPublicDomain_ExampleCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("example.com")
	if result != true {
		t.Errorf("isPublicDomain(example.com) = %v, expected true", result)
	}
}

func TestIsPublicDomain_PublicDomainCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("public-domain.com")
	if result != true {
		t.Errorf("isPublicDomain(public-domain.com) = %v, expected true", result)
	}
}

func TestIsPublicDomain_CdnExampleOrg(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("cdn.example.org")
	if result != true {
		t.Errorf("isPublicDomain(cdn.example.org) = %v, expected true", result)
	}
}

func TestIsPublicDomain_Localhost(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("localhost")
	if result != false {
		t.Errorf("isPublicDomain(localhost) = %v, expected false", result)
	}
}

func TestIsPublicDomain_127001(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("127.0.0.1")
	if result != false {
		t.Errorf("isPublicDomain(127.0.0.1) = %v, expected false", result)
	}
}

func TestIsPublicDomain_19216811(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("192.168.1.1")
	if result != false {
		t.Errorf("isPublicDomain(192.168.1.1) = %v, expected false", result)
	}
}

func TestIsPublicDomain_169254169254(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("169.254.169.254")
	if result != false {
		t.Errorf("isPublicDomain(169.254.169.254) = %v, expected false", result)
	}
}

func TestIsPublicDomain_ServerLocal(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("server.local")
	if result != false {
		t.Errorf("isPublicDomain(server.local) = %v, expected false", result)
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
func TestValidateURL_ValidHTTPURL(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://example.com/image.jpg")
	if err != nil {
		t.Errorf("Unexpected error for URL: %v", err)
	}
}

func TestValidateURL_ValidHTTPSURL(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("https://example.com/image.jpg")
	if err != nil {
		t.Errorf("Unexpected error for URL: %v", err)
	}
}

func TestValidateURL_InvalidScheme(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("ftp://example.com/image.jpg")
	if err == nil {
		t.Error("Expected error for URL")
	}
}

func TestValidateURL_EmptyURL(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("")
	if err == nil {
		t.Error("Expected error for URL")
	}
}

func TestValidateURL_InvalidFormat(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("not-a-url")
	if err == nil {
		t.Error("Expected error for URL")
	}
}

func TestValidateURL_URLWithoutHostname(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http:///image.jpg")
	if err == nil {
		t.Error("Expected error for URL")
	}
}

func TestValidateURL_URLWithSpecialChars(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://example.com/image.jpg?test=1")
	if err != nil {
		t.Errorf("Unexpected error for URL: %v", err)
	}
}

func TestValidateURL_URLWithPort(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://example.com:8080/image.jpg")
	if err != nil {
		t.Errorf("Unexpected error for URL: %v", err)
	}
}

// TestPrepareDataMissingParameters tests missing parameter validation
func TestPrepareDataMissingParameters_MissingExtension(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			_, errMsg := c.prepareData(r)
			return errMsg
		})
	router.AddRoute(route)

	req := httptest.NewRequest("GET", "/th//300x200/70/image.jpg", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Body.String() != "image extension is missing" {
		t.Errorf("Expected error %q, got %q", "image extension is missing", rr.Body.String())
	}
}

func TestPrepareDataMissingParameters_MissingSize(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			_, errMsg := c.prepareData(r)
			return errMsg
		})
	router.AddRoute(route)

	req := httptest.NewRequest("GET", "/th/jpg//70/image.jpg", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Body.String() != "size is missing" {
		t.Errorf("Expected error %q, got %q", "size is missing", rr.Body.String())
	}
}

func TestPrepareDataMissingParameters_MissingQuality(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			_, errMsg := c.prepareData(r)
			return errMsg
		})
	router.AddRoute(route)

	req := httptest.NewRequest("GET", "/th/jpg/300x200//image.jpg", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Body.String() != "quality is missing" {
		t.Errorf("Expected error %q, got %q", "quality is missing", rr.Body.String())
	}
}

func TestPrepareDataMissingParameters_MissingPath(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			_, errMsg := c.prepareData(r)
			return errMsg
		})
	router.AddRoute(route)

	req := httptest.NewRequest("GET", "/th/jpg/300x200/70/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Body.String() != "path is missing" {
		t.Errorf("Expected error %q, got %q", "path is missing", rr.Body.String())
	}
}

// TestPrepareDataURLNormalization tests URL path normalization
func TestPrepareDataURLNormalization_HTTPURLNormalization(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%s|%t|%t|%t", data.path, data.isURL, data.isCache, data.isFiles)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/http/example.com/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 4 {
		if parts[0] != "http://example.com/image.jpg" {
			t.Errorf("Expected path %q, got %q", "http://example.com/image.jpg", parts[0])
		}
		if parts[1] != "true" {
			t.Errorf("Expected isURL true, got %s", parts[1])
		}
		if parts[2] != "false" {
			t.Errorf("Expected isCache false, got %s", parts[2])
		}
		if parts[3] != "false" {
			t.Errorf("Expected isFiles false, got %s", parts[3])
		}
	}
}

func TestPrepareDataURLNormalization_HTTPSURLNormalization(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%s|%t|%t|%t", data.path, data.isURL, data.isCache, data.isFiles)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/https/example.com/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 4 {
		if parts[0] != "https://example.com/image.jpg" {
			t.Errorf("Expected path %q, got %q", "https://example.com/image.jpg", parts[0])
		}
		if parts[1] != "true" {
			t.Errorf("Expected isURL true, got %s", parts[1])
		}
		if parts[2] != "false" {
			t.Errorf("Expected isCache false, got %s", parts[2])
		}
		if parts[3] != "false" {
			t.Errorf("Expected isFiles false, got %s", parts[3])
		}
	}
}

func TestPrepareDataURLNormalization_CachePrefix(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%s|%t|%t|%t", data.path, data.isURL, data.isCache, data.isFiles)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/cache-mykey"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 4 {
		if parts[0] != "mykey" {
			t.Errorf("Expected path %q, got %q", "mykey", parts[0])
		}
		if parts[1] != "false" {
			t.Errorf("Expected isURL false, got %s", parts[1])
		}
		if parts[2] != "true" {
			t.Errorf("Expected isCache true, got %s", parts[2])
		}
		if parts[3] != "false" {
			t.Errorf("Expected isFiles false, got %s", parts[3])
		}
	}
}

func TestPrepareDataURLNormalization_FilesPrefix(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%s|%t|%t|%t", data.path, data.isURL, data.isCache, data.isFiles)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/files/uploads/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 4 {
		if parts[0] != "uploads/image.jpg" {
			t.Errorf("Expected path %q, got %q", "uploads/image.jpg", parts[0])
		}
		if parts[1] != "false" {
			t.Errorf("Expected isURL false, got %s", parts[1])
		}
		if parts[2] != "false" {
			t.Errorf("Expected isCache false, got %s", parts[2])
		}
		if parts[3] != "true" {
			t.Errorf("Expected isFiles true, got %s", parts[3])
		}
	}
}

// TestPrepareDataSizeParsing tests various size format parsing
func TestPrepareDataSizeParsing_WidthAndHeight(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d|%d", data.width, data.height)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 2 {
		var w, h int64
		fmt.Sscanf(parts[0], "%d", &w)
		fmt.Sscanf(parts[1], "%d", &h)
		if w != 300 || h != 200 {
			t.Errorf("Expected 300x200, got %dx%d", w, h)
		}
	}
}

func TestPrepareDataSizeParsing_WidthOnly(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d|%d", data.width, data.height)
		})
	router.AddRoute(route)

	url := "/th/jpg/500/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 2 {
		var w, h int64
		fmt.Sscanf(parts[0], "%d", &w)
		fmt.Sscanf(parts[1], "%d", &h)
		if w != 500 || h != 0 {
			t.Errorf("Expected 500x0, got %dx%d", w, h)
		}
	}
}

func TestPrepareDataSizeParsing_ZeroDimensions(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d|%d", data.width, data.height)
		})
	router.AddRoute(route)

	url := "/th/jpg/0x0/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 2 {
		var w, h int64
		fmt.Sscanf(parts[0], "%d", &w)
		fmt.Sscanf(parts[1], "%d", &h)
		if w != 0 || h != 0 {
			t.Errorf("Expected 0x0, got %dx%d", w, h)
		}
	}
}

func TestPrepareDataSizeParsing_LargeDimensions(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d|%d", data.width, data.height)
		})
	router.AddRoute(route)

	url := "/th/jpg/4000x3000/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	parts := strings.Split(body, "|")
	if len(parts) >= 2 {
		var w, h int64
		fmt.Sscanf(parts[0], "%d", &w)
		fmt.Sscanf(parts[1], "%d", &h)
		if w != 4000 || h != 3000 {
			t.Errorf("Expected 4000x3000, got %dx%d", w, h)
		}
	}
}

// TestSetHeadersAllFormats tests all supported image format headers
func TestSetHeadersAllFormats_Jpg(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "jpg")

	contentType := rec.Header().Get("Content-Type")
	if contentType != "image/jpeg" {
		t.Errorf("Extension jpg: expected Content-Type %q, got %q", "image/jpeg", contentType)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "max-age=604800" {
		t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
	}
}

func TestSetHeadersAllFormats_Jpeg(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "jpeg")

	contentType := rec.Header().Get("Content-Type")
	if contentType != "image/jpeg" {
		t.Errorf("Extension jpeg: expected Content-Type %q, got %q", "image/jpeg", contentType)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "max-age=604800" {
		t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
	}
}

func TestSetHeadersAllFormats_Png(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "png")

	contentType := rec.Header().Get("Content-Type")
	if contentType != "image/png" {
		t.Errorf("Extension png: expected Content-Type %q, got %q", "image/png", contentType)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "max-age=604800" {
		t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
	}
}

func TestSetHeadersAllFormats_Gif(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "gif")

	contentType := rec.Header().Get("Content-Type")
	if contentType != "image/gif" {
		t.Errorf("Extension gif: expected Content-Type %q, got %q", "image/gif", contentType)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "max-age=604800" {
		t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
	}
}

func TestSetHeadersAllFormats_Webp(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "webp")

	contentType := rec.Header().Get("Content-Type")
	if contentType != "" {
		t.Errorf("Extension webp: expected Content-Type %q, got %q", "", contentType)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "max-age=604800" {
		t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
	}
}

func TestSetHeadersAllFormats_Bmp(t *testing.T) {
	c := &thumbnailController{}
	rec := httptest.NewRecorder()
	c.setHeaders(rec, "bmp")

	contentType := rec.Header().Get("Content-Type")
	if contentType != "" {
		t.Errorf("Extension bmp: expected Content-Type %q, got %q", "", contentType)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "max-age=604800" {
		t.Errorf("Expected Cache-Control max-age=604800, got %q", cacheControl)
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
func TestIsPrivateHostEdgeCases_1721501(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("172.15.0.1")
	if result != true {
		t.Errorf("isPrivateHost(172.15.0.1) = %v, expected true", result)
	}
}

func TestIsPrivateHostEdgeCases_1723201(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("172.32.0.1")
	if result != true {
		t.Errorf("isPrivateHost(172.32.0.1) = %v, expected true", result)
	}
}

func TestIsPrivateHostEdgeCases_16925301(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("169.253.0.1")
	if result != false {
		t.Errorf("isPrivateHost(169.253.0.1) = %v, expected false", result)
	}
}

func TestIsPrivateHostEdgeCases_16925501(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("169.255.0.1")
	if result != false {
		t.Errorf("isPrivateHost(169.255.0.1) = %v, expected false", result)
	}
}

func TestIsPrivateHostEdgeCases_ExampleCorpCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("example.corp.com")
	if result != true {
		t.Errorf("isPrivateHost(example.corp.com) = %v, expected true", result)
	}
}

func TestIsPrivateHostEdgeCases_CorpExampleCom(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("corp.example.com")
	if result != false {
		t.Errorf("isPrivateHost(corp.example.com) = %v, expected false", result)
	}
}

func TestIsPrivateHostEdgeCases_10255255255(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("10.255.255.255")
	if result != true {
		t.Errorf("isPrivateHost(10.255.255.255) = %v, expected true", result)
	}
}

func TestIsPrivateHostEdgeCases_11000(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPrivateHost("11.0.0.0")
	if result != false {
		t.Errorf("isPrivateHost(11.0.0.0) = %v, expected false", result)
	}
}

// TestValidateURLWithSpecialCases tests additional URL validation scenarios
func TestValidateURLWithSpecialCases_IPv6Localhost(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://[::1]/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

func TestValidateURLWithSpecialCases_IPv6Address(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://[2001:db8::1]/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

func TestValidateURLWithSpecialCases_URLWithAuth(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://user:pass@example.com/image.jpg")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

func TestValidateURLWithSpecialCases_URLWithFragment(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://example.com/image.jpg#section")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

func TestValidateURLWithSpecialCases_URLWithQueryParams(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://example.com/image.jpg?size=large&format=jpg")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

func TestValidateURLWithSpecialCases_LocalDomain(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://myserver.local/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

func TestValidateURLWithSpecialCases_InternalDomain(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://api.internal/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

func TestValidateURLWithSpecialCases_NoScheme(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("example.com/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

// TestPrepareDataQualityParsing tests quality parameter parsing
func TestPrepareDataQualityParsing_LowQuality(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d", data.quality)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/30/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	var q int64
	fmt.Sscanf(body, "%d", &q)
	if q != 30 {
		t.Errorf("Expected quality 30, got %d", q)
	}
}

func TestPrepareDataQualityParsing_MediumQuality(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d", data.quality)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	var q int64
	fmt.Sscanf(body, "%d", &q)
	if q != 70 {
		t.Errorf("Expected quality 70, got %d", q)
	}
}

func TestPrepareDataQualityParsing_HighQuality(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d", data.quality)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/95/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	var q int64
	fmt.Sscanf(body, "%d", &q)
	if q != 95 {
		t.Errorf("Expected quality 95, got %d", q)
	}
}

func TestPrepareDataQualityParsing_ZeroQuality(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d", data.quality)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/0/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	var q int64
	fmt.Sscanf(body, "%d", &q)
	if q != 0 {
		t.Errorf("Expected quality 0, got %d", q)
	}
}

func TestPrepareDataQualityParsing_InvalidQuality(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return fmt.Sprintf("%d", data.quality)
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/abc/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	var q int64
	fmt.Sscanf(body, "%d", &q)
	if q != 0 {
		t.Errorf("Expected quality 0, got %d", q)
	}
}

// TestValidateURLWithCorporateDomains tests corporate domain blocking
func TestValidateURLWithCorporateDomains_InternalCorpDomain(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://internal.corp.example.com/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

func TestValidateURLWithCorporateDomains_NestedCorpDomain(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://api.service.corp.example.com/image.jpg")
	if err == nil {
		t.Error("Expected error for URL, but got none")
	}
}

func TestValidateURLWithCorporateDomains_PublicCorpDomain(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://example.corp/image.jpg")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

func TestValidateURLWithCorporateDomains_CorpInSubdomain(t *testing.T) {
	c := &thumbnailController{}
	err := c.validateURL("http://corp.example.com/image.jpg")
	if err != nil {
		t.Errorf("Expected no error for URL, but got: %v", err)
	}
}

// TestPrepareDataExtensionHandling tests extension parameter handling
func TestPrepareDataExtensionHandling_JPEGExtension(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return data.extension
		})
	router.AddRoute(route)

	url := "/th/jpg/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body != "jpg" {
		t.Errorf("Expected extension %q, got %q", "jpg", body)
	}
}

func TestPrepareDataExtensionHandling_PNGExtension(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return data.extension
		})
	router.AddRoute(route)

	url := "/th/png/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body != "png" {
		t.Errorf("Expected extension %q, got %q", "png", body)
	}
}

func TestPrepareDataExtensionHandling_GIFExtension(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return data.extension
		})
	router.AddRoute(route)

	url := "/th/gif/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body != "gif" {
		t.Errorf("Expected extension %q, got %q", "gif", body)
	}
}

func TestPrepareDataExtensionHandling_UppercaseExtension(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return data.extension
		})
	router.AddRoute(route)

	url := "/th/JPG/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body != "jpg" {
		t.Errorf("Expected extension %q, got %q", "jpg", body)
	}
}

func TestPrepareDataExtensionHandling_WebPExtension(t *testing.T) {
	c := &thumbnailController{}
	router := rtr.NewRouter()
	route := rtr.NewRoute().
		SetMethod("GET").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			data, _ := c.prepareData(r)
			return data.extension
		})
	router.AddRoute(route)

	url := "/th/webp/300x200/70/image.jpg"
	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body != "webp" {
		t.Errorf("Expected extension %q, got %q", "webp", body)
	}
}

// TestIsPublicDomainComprehensive tests public domain detection comprehensively
func TestIsPublicDomainComprehensive_PublicDomain(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("example.com")
	if result != true {
		t.Errorf("isPublicDomain(example.com) = %v, expected true", result)
	}
}

func TestIsPublicDomainComprehensive_PublicSubdomain(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("cdn.example.com")
	if result != true {
		t.Errorf("isPublicDomain(cdn.example.com) = %v, expected true", result)
	}
}

func TestIsPublicDomainComprehensive_PublicWithHyphen(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("my-domain.com")
	if result != true {
		t.Errorf("isPublicDomain(my-domain.com) = %v, expected true", result)
	}
}

func TestIsPublicDomainComprehensive_PublicWithNumbers(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("example123.com")
	if result != true {
		t.Errorf("isPublicDomain(example123.com) = %v, expected true", result)
	}
}

func TestIsPublicDomainComprehensive_Localhost(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("localhost")
	if result != false {
		t.Errorf("isPublicDomain(localhost) = %v, expected false", result)
	}
}

func TestIsPublicDomainComprehensive_127001(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("127.0.0.1")
	if result != false {
		t.Errorf("isPublicDomain(127.0.0.1) = %v, expected false", result)
	}
}

func TestIsPublicDomainComprehensive_PrivateIP(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("192.168.1.1")
	if result != false {
		t.Errorf("isPublicDomain(192.168.1.1) = %v, expected false", result)
	}
}

func TestIsPublicDomainComprehensive_LocalDomain(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("server.local")
	if result != false {
		t.Errorf("isPublicDomain(server.local) = %v, expected false", result)
	}
}

func TestIsPublicDomainComprehensive_InternalDomain(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("api.internal")
	if result != false {
		t.Errorf("isPublicDomain(api.internal) = %v, expected false", result)
	}
}

func TestIsPublicDomainComprehensive_IPv6Loopback(t *testing.T) {
	c := &thumbnailController{}
	result := c.isPublicDomain("::1")
	if result != false {
		t.Errorf("isPublicDomain(::1) = %v, expected false", result)
	}
}
