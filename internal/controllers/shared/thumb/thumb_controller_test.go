package thumb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

// helper to build a request with chi URL params used by prepareData
func newChiRequest(method, target string, params map[string]string) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
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

// TestPrepareDataValidations checks missing params produce clear errors
func TestPrepareDataValidations(t *testing.T) {
	c := &thumbnailController{}

	// missing extension
	r := newChiRequest(http.MethodGet, "/th//200x100/80/path.jpg", map[string]string{
		"extension": "",
		"size":      "200x100",
		"quality":   "80",
		"*":         "path.jpg",
	})
	_, errMsg := c.prepareData(r)
	if errMsg != "image extension is missing" {
		t.Fatalf("expected 'image extension is missing', got: %q", errMsg)
	}

	// missing size
	r = newChiRequest(http.MethodGet, "/th/jpg//80/path.jpg", map[string]string{
		"extension": "jpg",
		"size":      "",
		"quality":   "80",
		"*":         "path.jpg",
	})
	_, errMsg = c.prepareData(r)
	if errMsg != "size is missing" {
		t.Fatalf("expected 'size is missing', got: %q", errMsg)
	}

	// missing quality
	r = newChiRequest(http.MethodGet, "/th/jpg/200x100//path.jpg", map[string]string{
		"extension": "jpg",
		"size":      "200x100",
		"quality":   "",
		"*":         "path.jpg",
	})
	_, errMsg = c.prepareData(r)
	if errMsg != "quality is missing" {
		t.Fatalf("expected 'quality is missing', got: %q", errMsg)
	}

	// missing path
	r = newChiRequest(http.MethodGet, "/th/jpg/200x100/80/", map[string]string{
		"extension": "jpg",
		"size":      "200x100",
		"quality":   "80",
		"*":         "",
	})
	_, errMsg = c.prepareData(r)
	if errMsg != "path is missing" {
		t.Fatalf("expected 'path is missing', got: %q", errMsg)
	}
}

// TestPrepareDataParsing verifies size parsing and URL normalization
func TestPrepareDataParsing(t *testing.T) {
	c := &thumbnailController{}

	// width x height parsing
	r := newChiRequest(http.MethodGet, "/th/jpg/200x150/70/a/b.jpg", map[string]string{
		"extension": "jpg",
		"size":      "200x150",
		"quality":   "70",
		"*":         "a/b.jpg",
	})
	data, errMsg := c.prepareData(r)
	if errMsg != "" {
		t.Fatalf("unexpected error: %q", errMsg)
	}
	if data.width != 200 || data.height != 150 || data.quality != 70 {
		t.Fatalf("expected width=200 height=150 quality=70, got %d %d %d", data.width, data.height, data.quality)
	}

	// single dimension parsing
	r = newChiRequest(http.MethodGet, "/th/png/300/60/a/b.png", map[string]string{
		"extension": "png",
		"size":      "300",
		"quality":   "60",
		"*":         "a/b.png",
	})
	data, errMsg = c.prepareData(r)
	if errMsg != "" {
		t.Fatalf("unexpected error: %q", errMsg)
	}
	if data.width != 300 || data.height != 0 {
		t.Fatalf("expected width=300 height=0, got %d %d", data.width, data.height)
	}

	// URL normalization (https prefix)
	r = newChiRequest(http.MethodGet, "/th/jpg/100x0/70/https/example.com/img.jpg", map[string]string{
		"extension": "jpg",
		"size":      "100x0",
		"quality":   "70",
		"*":         "https/example.com/img.jpg",
	})
	data, errMsg = c.prepareData(r)
	if errMsg != "" {
		t.Fatalf("unexpected error: %q", errMsg)
	}
	if !data.isURL || data.path != "https://example.com/img.jpg" {
		t.Fatalf("expected isURL=true and normalized https URL, got isURL=%v path=%q", data.isURL, data.path)
	}
}

// TestPrepareDataFilesLink ensures files/ paths are converted via links.URL and marked as URL
func TestPrepareDataFilesLink(t *testing.T) {
	// Ensure APP_ENV=testing so links.RootURL() returns empty string (see links.RootURL)
	oldEnv := os.Getenv("APP_ENV")
	_ = os.Setenv("APP_ENV", "testing")
	defer os.Setenv("APP_ENV", oldEnv)

	c := &thumbnailController{}
	r := newChiRequest(http.MethodGet, "/th/jpg/100x0/70/files/a/b.jpg", map[string]string{
		"extension": "jpg",
		"size":      "100x0",
		"quality":   "70",
		"*":         "files/a/b.jpg",
	})
	data, errMsg := c.prepareData(r)
	if errMsg != "" {
		t.Fatalf("unexpected error: %q", errMsg)
	}
	if !data.isURL {
		t.Fatalf("expected isURL=true for files/ path")
	}
	// With APP_ENV=testing, links.URL returns "/" + path
	if data.path != "/files/a/b.jpg" {
		t.Fatalf("expected path '/files/a/b.jpg', got %q", data.path)
	}
}
