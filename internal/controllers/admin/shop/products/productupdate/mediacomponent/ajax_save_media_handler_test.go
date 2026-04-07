package mediacomponent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestHandleAjaxSaveMedia_Success tests successful media saving
func TestHandleAjaxSaveMedia_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("100.00")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	body := `{"media": [{"id": "1", "fileName": "test.jpg", "url": "https://example.com/test.jpg", "isMain": true}]}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-media&product_id="+product.GetID(),
		strings.NewReader(body))

	response := HandleAjaxSaveMedia(registry, req, product.GetID())

	if !strings.Contains(response, `"media"`) {
		t.Error("expected response to contain media field")
	}
	if !strings.Contains(response, `"test.jpg"`) {
		t.Error("expected response to contain test.jpg")
	}

	// Verify media was actually saved
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(product.GetID())
	medias, err := registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	if err != nil {
		t.Fatalf("failed to list media: %v", err)
	}

	if len(medias) != 1 {
		t.Errorf("expected 1 media, got %d", len(medias))
	}

	if len(medias) > 0 && medias[0].GetURL() != "https://example.com/test.jpg" {
		t.Errorf("expected URL https://example.com/test.jpg, got %s", medias[0].GetURL())
	}
	if len(medias) > 0 && medias[0].GetType() != "image" {
		t.Errorf("expected type 'image', got '%s'", medias[0].GetType())
	}
	if len(medias) > 0 && medias[0].GetSequence() != 0 {
		t.Errorf("expected sequence 0, got %d", medias[0].GetSequence())
	}
}

// TestHandleAjaxSaveMedia_MultipleItems tests saving multiple media with correct sequence
func TestHandleAjaxSaveMedia_MultipleItems(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("100.00")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	body := `{"media": [
		{"id": "1", "fileName": "first.jpg", "url": "https://example.com/first.jpg", "isMain": true},
		{"id": "2", "fileName": "second.mp4", "url": "https://example.com/second.mp4", "isMain": false},
		{"id": "3", "fileName": "third.png", "url": "https://example.com/third.png", "isMain": false}
	]}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-media&product_id="+product.GetID(),
		strings.NewReader(body))

	response := HandleAjaxSaveMedia(registry, req, product.GetID())

	if !strings.Contains(response, `"media"`) {
		t.Error("expected response to contain media field")
	}

	// Verify media was saved with correct sequence
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(product.GetID())
	medias, err := registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	if err != nil {
		t.Fatalf("failed to list media: %v", err)
	}

	if len(medias) != 3 {
		t.Errorf("expected 3 media, got %d", len(medias))
	}

	// Check types
	for _, media := range medias {
		switch media.GetURL() {
		case "https://example.com/first.jpg", "https://example.com/third.png":
			if media.GetType() != "image" {
				t.Errorf("expected type 'image' for %s, got '%s'", media.GetURL(), media.GetType())
			}
		case "https://example.com/second.mp4":
			if media.GetType() != "video" {
				t.Errorf("expected type 'video' for %s, got '%s'", media.GetURL(), media.GetType())
			}
		}
	}
}

// TestDetermineMediaType tests the media type determination function
func TestDetermineMediaType(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/photo.jpg", "image"},
		{"https://example.com/photo.jpeg", "image"},
		{"https://example.com/photo.png", "image"},
		{"https://example.com/photo.gif", "image"},
		{"https://example.com/video.mp4", "video"},
		{"https://example.com/video.webm", "video"},
		{"https://example.com/audio.mp3", "audio"},
		{"https://example.com/audio.wav", "audio"},
		{"https://example.com/document.pdf", "file"},
		{"https://example.com/unknown.xyz", "file"},
	}

	for _, tt := range tests {
		result := determineMediaType(tt.url)
		if result != tt.expected {
			t.Errorf("determineMediaType(%q) = %q, expected %q", tt.url, result, tt.expected)
		}
	}
}
func TestHandleAjaxSaveMedia_InvalidBody(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-media&product_id=test123",
		strings.NewReader(body))

	response := HandleAjaxSaveMedia(registry, req, "test123")

	if !strings.Contains(response, `"error"`) {
		t.Error("expected response to contain error field")
	}
	if !strings.Contains(response, `Invalid request body`) {
		t.Error("expected response to contain 'Invalid request body' message")
	}
}

// TestHandleAjaxSaveMedia_EmptyMedia tests saving empty media list
func TestHandleAjaxSaveMedia_EmptyMedia(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	product := shopstore.NewProduct()
	product.SetTitle("Test Product")
	product.SetPrice("100.00")
	product.SetQuantity("10")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)

	if err := registry.GetShopStore().ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// Create existing media
	media := shopstore.NewMedia()
	media.SetEntityID(product.GetID())
	media.SetURL("https://example.com/old.jpg")
	media.SetTitle("Old")
	media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	media.SetSequence(0)
	media.SetType("image")
	if err := registry.GetShopStore().MediaCreate(context.Background(), media); err != nil {
		t.Fatalf("failed to create media: %v", err)
	}

	// Save empty list - should delete existing
	body := `{"media": []}`
	req := httptest.NewRequest(http.MethodPost, "/?action=save-media&product_id="+product.GetID(),
		strings.NewReader(body))

	response := HandleAjaxSaveMedia(registry, req, product.GetID())

	if !strings.Contains(response, `"media"`) {
		t.Error("expected response to contain media field")
	}

	// Verify old media was deleted
	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(product.GetID())
	medias, err := registry.GetShopStore().MediaList(context.Background(), mediaQuery)
	if err != nil {
		t.Fatalf("failed to list media: %v", err)
	}

	if len(medias) != 0 {
		t.Errorf("expected 0 media after empty save, got %d", len(medias))
	}
}
