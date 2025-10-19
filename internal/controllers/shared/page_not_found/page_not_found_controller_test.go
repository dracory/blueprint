package page_not_found

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandlerReturnsNotFound ensures the controller sets 404 and returns message
func TestHandlerReturnsNotFound(t *testing.T) {
	c := PageNotFoundController()

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	res := c.Handler(rec, req)

	// Verify status code written to ResponseWriter
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	// Verify HTML content contains key elements
	if !strings.Contains(res, "<title>404 - Page Not Found</title>") {
		t.Fatal("expected HTML response with title")
	}
	if !strings.Contains(res, "Oops! Page Not Found") {
		t.Fatal("expected HTML response with error title")
	}
}
