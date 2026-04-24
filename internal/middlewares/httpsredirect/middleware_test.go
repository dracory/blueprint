package httpsredirect

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func serve(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	NewHTTPSRedirectMiddleware().GetHandler()(okHandler).ServeHTTP(rr, r)
	return rr
}

func newRequest(url string) *http.Request {
	return httptest.NewRequest(http.MethodGet, url, nil)
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Errorf("got status %v, want %v", rr.Code, want)
	}
}

func assertLocation(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()
	if loc := rr.Header().Get("Location"); loc != want {
		t.Errorf("got Location %q, want %q", loc, want)
	}
}

func TestRedirectsHTTPToHTTPS(t *testing.T) {
	req := newRequest("http://example.com/foo?bar=baz")
	rr := serve(req)
	assertStatus(t, rr, http.StatusMovedPermanently)
	assertLocation(t, rr, "https://example.com/foo?bar=baz")
}

func TestPreservesPathAndQueryOnRedirect(t *testing.T) {
	req := newRequest("http://example.com/some/path?foo=bar&baz=qux")
	rr := serve(req)
	assertStatus(t, rr, http.StatusMovedPermanently)
	assertLocation(t, rr, "https://example.com/some/path?foo=bar&baz=qux")
}

func TestNoRedirectForHTTPSViaTLS(t *testing.T) {
	req := newRequest("http://example.com/foo")
	req.TLS = &tls.ConnectionState{}
	assertStatus(t, serve(req), http.StatusOK)
}

func TestNoRedirectForHTTPSViaXForwardedProto(t *testing.T) {
	req := newRequest("http://example.com/foo")
	req.Header.Set("X-Forwarded-Proto", "https")
	assertStatus(t, serve(req), http.StatusOK)
}

func TestNoRedirectForHTTPSViaXForwardedScheme(t *testing.T) {
	req := newRequest("http://example.com/foo")
	req.Header.Set("X-Forwarded-Scheme", "https")
	assertStatus(t, serve(req), http.StatusOK)
}

func TestNoRedirectForLocalhost(t *testing.T) {
	assertStatus(t, serve(newRequest("http://localhost/foo")), http.StatusOK)
}

func TestNoRedirectForLocalhostWithPort(t *testing.T) {
	assertStatus(t, serve(newRequest("http://localhost:8080/foo")), http.StatusOK)
}

func TestNoRedirectFor127(t *testing.T) {
	assertStatus(t, serve(newRequest("http://127.0.0.1/foo")), http.StatusOK)
}

func TestNoRedirectForIPv6Loopback(t *testing.T) {
	assertStatus(t, serve(newRequest("http://[::1]/foo")), http.StatusOK)
}

func TestNoRedirectFor10xPrivateRange(t *testing.T) {
	assertStatus(t, serve(newRequest("http://10.0.0.1/foo")), http.StatusOK)
}

func TestNoRedirectFor192168PrivateRange(t *testing.T) {
	assertStatus(t, serve(newRequest("http://192.168.1.1/foo")), http.StatusOK)
}

func TestNoRedirectForLocalDomain(t *testing.T) {
	assertStatus(t, serve(newRequest("http://myapp.local/foo")), http.StatusOK)
}

func TestNoRedirectInDevelopmentEnv(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	assertStatus(t, serve(newRequest("http://example.com/foo")), http.StatusOK)
}

func TestNoRedirectInLocalEnv(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	assertStatus(t, serve(newRequest("http://example.com/foo")), http.StatusOK)
}

func TestNoRedirectInTestingEnv(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	assertStatus(t, serve(newRequest("http://example.com/foo")), http.StatusOK)
}
