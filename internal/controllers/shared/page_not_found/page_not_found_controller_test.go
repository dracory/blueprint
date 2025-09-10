package page_not_found

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerReturnsNotFound ensures the controller sets 404 and returns message
func TestHandlerReturnsNotFound(t *testing.T) {
	c := PageNotFoundController()

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	res := c.Handler(rec, req)

	// Verify returned message
	if res != "Sorry, page not found." {
		t.Fatalf("expected 'Sorry, page not found.', got: %q", res)
	}

	// Verify status code written to ResponseWriter
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	// Handler does not write body, only status; ensure body is empty
	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty response body, got: %q", rec.Body.String())
	}
}
