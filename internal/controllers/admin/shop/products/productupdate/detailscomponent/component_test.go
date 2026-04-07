package detailscomponent

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

	if !strings.Contains(html, "details-wrapper") {
		t.Error("expected HTML to contain details-wrapper div")
	}
	if !strings.Contains(html, "vue.global.js") {
		t.Error("expected HTML to contain Vue CDN")
	}
	if !strings.Contains(html, "quill") {
		t.Error("expected HTML to contain Quill editor references")
	}
	if !strings.Contains(html, productID) {
		t.Error("expected HTML to contain product ID")
	}
}

// TestRender_ContainsRequiredScripts tests that required scripts are included
func TestRender_ContainsRequiredScripts(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	productID := "test-product-456"

	result := Render(registry, productID)
	html := result.ToHTML()

	// Check for Vue.js
	if !strings.Contains(html, "vue@3") {
		t.Error("expected HTML to contain Vue.js CDN")
	}

	// Check for Quill
	if !strings.Contains(html, "quill@1.3.7") {
		t.Error("expected HTML to contain Quill CDN")
	}

	// Check for Vue Quill
	if !strings.Contains(html, "@vueup/vue-quill") {
		t.Error("expected HTML to contain Vue Quill CDN")
	}
}

// TestRender_ContainsFormWrapper tests that form wrapper is present
func TestRender_ContainsFormWrapper(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithShopStore(true),
	)

	productID := "test-product-789"

	result := Render(registry, productID)
	html := result.ToHTML()

	if !strings.Contains(html, "FormProductDetailsUpdate") {
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
	if !strings.Contains(html, "urlDetailsLoad") {
		t.Error("expected HTML to contain urlDetailsLoad variable")
	}
	if !strings.Contains(html, "urlDetailsSave") {
		t.Error("expected HTML to contain urlDetailsSave variable")
	}
}
