package metadatacomponent

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

// TestRender_Success tests successful component rendering
func TestRender_Success(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	productID := "test-product-123"

	result := Render(registry, productID)

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	html := result.ToHTML()

	if !strings.Contains(html, "metadata-wrapper") {
		t.Error("expected HTML to contain metadata-wrapper div")
	}
	if !strings.Contains(html, productID) {
		t.Error("expected HTML to contain product ID")
	}
}

// TestRender_ContainsVueJS tests that Vue.js is included
func TestRender_ContainsVueJS(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	productID := "test-product-vue"

	result := Render(registry, productID)
	html := result.ToHTML()

	if !strings.Contains(html, "vue.global.js") || !strings.Contains(html, "vue@3") {
		t.Error("expected HTML to contain Vue.js CDN")
	}
}

// TestRender_ContainsFormWrapper tests that form wrapper is present
func TestRender_ContainsFormWrapper(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	productID := "test-product-form"

	result := Render(registry, productID)
	html := result.ToHTML()

	if !strings.Contains(html, "FormProductMetadataUpdate") {
		t.Error("expected HTML to contain form wrapper ID")
	}
}

// TestRender_ContainsInitScript tests that initialization script is present
func TestRender_ContainsInitScript(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	productID := "test-product-init"

	result := Render(registry, productID)
	html := result.ToHTML()

	if !strings.Contains(html, "productId") {
		t.Error("expected HTML to contain productId variable")
	}
	if !strings.Contains(html, "urlMetasLoad") {
		t.Error("expected HTML to contain urlMetasLoad variable")
	}
	if !strings.Contains(html, "urlMetasSave") {
		t.Error("expected HTML to contain urlMetasSave variable")
	}
}
