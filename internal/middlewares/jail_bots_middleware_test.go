package middlewares

import (
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/testutils"

	"github.com/jellydator/ttlcache/v3"
	"github.com/spf13/cast"
)

func TestJailBotsMiddlewareName(t *testing.T) {

	// Act
	m := JailBotsMiddleware(JailBotsConfig{})

	// Assert
	if m.GetName() != "Jail Bots Middleware" {
		t.Fatal("JailBotsMiddleware.Name must be Jail Bots Middleware. Got ", m.GetName())
	}
}

func TestJailBotsMiddlewareAllowedResponse(t *testing.T) {

	allowedUris := []string{
		"/robots.txt",
		"/sitemap.xml",
		"/favicon.ico",
		"/",
		"/auth/login",
	}

	// Act

	for _, allowedUri := range allowedUris {
		m := JailBotsMiddleware(JailBotsConfig{})
		req, err := testutils.NewRequest("GET", allowedUri, testutils.NewRequestOptions{})

		if err != nil {
			t.Fatal(err)
		}

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// DEBUG: t.Log("Passes as expected")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("gone through"))
		})

		rw := httptest.NewRecorder()
		handler := m.Execute(testHandler)
		handler.ServeHTTP(rw, req)

		// Assert
		body := rw.Body.String()

		if body != "gone through" {
			t.Fatal("Response SHOULD NOT BE 'gone through' but found: ", rw.Body.String())
		}
	}
}

func TestJailBotsMiddlewareJailedResponse(t *testing.T) {

	allowedUris := []string{
		"/.env",
		"/backup",
		"/db",
		"/wp-admin",
	}

	// Act

	for _, allowedUri := range allowedUris {
		m := JailBotsMiddleware(JailBotsConfig{})
		randInt := rand.IntN(1000)
		req, err := testutils.NewRequest("GET", allowedUri, testutils.NewRequestOptions{
			Headers: map[string]string{
				"User-Agent":      "test-agent",
				"Referer":         "test-referer",
				"X-Forwarded-For": "127.0.0." + cast.ToString(randInt),
			},
		})

		if err != nil {
			t.Fatal(err)
		}

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// DEBUG: t.Log("Passes as expected")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("gone through"))
		})

		rw := httptest.NewRecorder()
		handler := m.Execute(testHandler)
		handler.ServeHTTP(rw, req)

		// Assert
		body := rw.Body.String()

		if body != "malicious access not allowed (jb)" {
			t.Fatal("Response SHOULD NOT BE 'malicious access not allowed (jb)' but found: ", rw.Body.String())
		}
	}
}

func TestJailBotsMiddlewareHandler(t *testing.T) {

	// Act
	m := JailBotsMiddleware(JailBotsConfig{})

	// Assert
	if m.GetHandler() == nil {
		t.Error("JailBotsMiddleware.Handler is nil")
	}
}

func TestJailBotsMiddlewareIsJailable(t *testing.T) {

	data := []struct {
		uri      string
		jailable bool
	}{
		{uri: "/", jailable: false},
		{uri: "/robots.txt", jailable: false},
		{uri: "/sitemap.xml", jailable: false},
		{uri: "/favicon.ico", jailable: false},
		{uri: "/.env", jailable: true},
		{uri: "/.well-known/ALFA_DATA", jailable: true},
		{uri: "/.well-known/alfacgiapi", jailable: true},
		{uri: "/.well-known/cgialfa", jailable: true},
		{uri: "/api/search?folderIds=0", jailable: true},
		{uri: "/aws/credentials", jailable: true},
		{uri: "/backup", jailable: true},
		{uri: "/backup/license.txt", jailable: true},
		{uri: "/bc", jailable: true},
		{uri: "/bk", jailable: true},
		{uri: "/blog/license.txt", jailable: true},
		{uri: "/bin/", jailable: true},
		{uri: "/cgialfa", jailable: true},
		{uri: "/content/sitetree", jailable: true},
		{uri: "/config.json", jailable: true},
		{uri: "/cgi-bin", jailable: true},
		{uri: "/credentials", jailable: true},
		{uri: "/db", jailable: true},
		{uri: "/db/license.txt", jailable: true},
		{uri: "/wp-config.php", jailable: true},
		{uri: "/wp-content/plugins", jailable: true},
		{uri: "/wp-content/themes", jailable: true},
		{uri: "/wp-includes", jailable: true},
		{uri: "/wp-includes/css", jailable: true},
		{uri: "/wp-includes/js", jailable: true},
		{uri: "/wp-includes/images", jailable: true},
		{uri: "/wp-includes/javascript", jailable: true},
	}

	// Act
	m := jailBotsMiddleware{}

	// Assert
	for i := 0; i < len(data); i++ {
		jailable, _ := m.isJailable(data[i].uri)
		if jailable != data[i].jailable {
			t.Fatal("JailBotsMiddleware.isJailable(", data[i].uri, ") must be ", data[i].jailable)
		}
	}
}

// helper to create a simple next handler that writes a body for assertion
func nextOK() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("next"))
	})
}

func TestExcludedWildcard_Blog(t *testing.T) {
	jb := &jailBotsMiddleware{
		cache:        ttlcache.New[string, struct{}](),
		excludePaths: []string{"/blog*"},
	}

	handler := jb.Handler(nextOK())

	req := httptest.NewRequest(http.MethodGet, "/blog/post", nil)
	req.RemoteAddr = "1.2.3.4:12345"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 for excluded wildcard path, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "next" {
		t.Fatalf("expected body 'next', got %q", string(body))
	}
}

func TestExcludedSegment_NoWildcard_AllowsBlogAndDoesNotOvermatch(t *testing.T) {
	jb := &jailBotsMiddleware{
		cache:        ttlcache.New[string, struct{}](),
		excludePaths: []string{"/blog"},
	}

	handler := jb.Handler(nextOK())

	// Exact or segment '/blog/...'
	req1 := httptest.NewRequest(http.MethodGet, "/blog/.env", nil)
	req1.RemoteAddr = "2.2.2.2:12345"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("expected 200 for excluded segment path '/blog/.env', got %d", rr1.Code)
	}

	// Ensure no overmatch like '/blogger'
	req2 := httptest.NewRequest(http.MethodGet, "/blogger", nil)
	req2.RemoteAddr = "2.2.2.3:12345"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200 pass-through for '/blogger', got %d", rr2.Code)
	}
}

func TestJailableBlocksAndJailsIP(t *testing.T) {
	jb := &jailBotsMiddleware{
		cache:        ttlcache.New[string, struct{}](),
		excludePaths: []string{},
	}

	handler := jb.Handler(nextOK())

	// First: trigger jailing via a known blacklisted prefix '/wp'
	req1 := httptest.NewRequest(http.MethodGet, "/wp-admin", nil)
	req1.RemoteAddr = "9.9.9.9:1111"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for jailable path, got %d", rr1.Code)
	}

	// Second: same IP to a normal path should still be forbidden due to jail
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "9.9.9.9:9999"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for jailed IP on normal path, got %d", rr2.Code)
	}
}

func TestExcludedOverridesJail(t *testing.T) {
	jb := &jailBotsMiddleware{
		cache:        ttlcache.New[string, struct{}](),
		excludePaths: []string{"/blog*"},
	}

	handler := jb.Handler(nextOK())

	// Jail the IP first
	req1 := httptest.NewRequest(http.MethodGet, "/wp-admin", nil)
	req1.RemoteAddr = "3.3.3.3:2222"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusForbidden {
		t.Fatalf("expected 403 on jailable path, got %d", rr1.Code)
	}

	// Now access excluded path with same IP should bypass jail and pass
	req2 := httptest.NewRequest(http.MethodGet, "/blog/post", nil)
	req2.RemoteAddr = "3.3.3.3:3333"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200 on excluded path even when jailed, got %d", rr2.Code)
	}
}
