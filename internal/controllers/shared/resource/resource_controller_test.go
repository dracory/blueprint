package resource

import (
	"net/http"
	"net/http/httptest"
	"project/internal/controllers/shared/page_not_found"
	"project/internal/resources"
	"testing"
)

// TestHandlerServesEmbeddedJS ensures an existing embedded resource is served
// with the correct Content-Type.
func TestHandlerServesEmbeddedJS(t *testing.T) {
	c := NewResourceController()

	req := httptest.NewRequest(http.MethodGet, "/resources/js/blockarea_v0200.js", nil)
	rec := httptest.NewRecorder()

	res := c.Handler(rec, req)

	expected, err := resources.ToString("js/blockarea_v0200.js")
	if err != nil {
		t.Fatalf("failed to read embedded resource for test: %v", err)
	}

	if res != expected {
		t.Fatalf("expected body to equal embedded resource content (%d bytes), got different content (%d bytes)", len(expected), len(res))
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/javascript" {
		t.Fatalf("expected Content-Type 'application/javascript', got: %q", ct)
	}
}

// TestHandlerPrivateResource ensures paths starting with '.' return 404 and PNf message.
func TestHandlerPrivateResource(t *testing.T) {
	c := NewResourceController()

	req := httptest.NewRequest(http.MethodGet, "/resources/.secrets.txt", nil)
	rec := httptest.NewRecorder()

	res := c.Handler(rec, req)

	if res != page_not_found.PageNotFoundController().Handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)) {
		t.Fatalf("expected page-not-found message, got: %q", res)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// TestHandlerMissingResource ensures missing resources return 404 via page-not-found.
func TestHandlerMissingResource(t *testing.T) {
	c := NewResourceController()

	req := httptest.NewRequest(http.MethodGet, "/resources/js/does_not_exist_123456.js", nil)
	rec := httptest.NewRecorder()

	res := c.Handler(rec, req)

	if res != page_not_found.PageNotFoundController().Handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)) {
		t.Fatalf("expected page-not-found message, got: %q", res)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
