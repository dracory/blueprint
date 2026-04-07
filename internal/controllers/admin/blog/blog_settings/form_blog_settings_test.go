package blog_settings

import (
	"context"
	"html"
	"strings"
	"testing"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/testutils"
)

func TestBlogSettingsForm_RenderStructure(t *testing.T) {
	// Create a mock app for testing
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	component := NewFormBlogSettings(registry)
	if component == nil {
		t.Fatal("Failed to create component")
	}

	// Set test data on the component
	if comp, ok := component.(*formBlogSettings); ok {
		comp.BlogTopic = "AI insights"
		comp.ReturnURL = shared.NewLinks().PostManager()
	}

	html := component.Render(context.TODO()).ToHTML()
	decoded := htmlDecode(html)
	t.Log(decoded)

	if !strings.Contains(html, `id="BlogSettingsFormWrapper"`) {
		t.Error("expected BlogSettingsFormWrapper id in HTML")
	}
	if !strings.Contains(html, `name="blog_topic"`) {
		t.Error("expected blog_topic name in HTML")
	}
	if !strings.Contains(html, "AI insights") {
		t.Error("expected AI insights in HTML")
	}
	cancelURL := shared.NewLinks().PostManager()
	if !strings.Contains(html, cancelURL) {
		t.Error("expected cancel URL in HTML")
	}

	if !strings.Contains(html, `data-flux-action="apply"`) {
		t.Error("expected apply action in HTML")
	}
	if !strings.Contains(html, `data-flux-action="save_close"`) {
		t.Error("expected save_close action in HTML")
	}

	// Verify the action attributes are present in the HTML
	if !strings.Contains(decoded, "data-flux-action=\"apply\"") {
		t.Error("expected apply action in decoded HTML")
	}
	if !strings.Contains(decoded, "data-flux-action=\"save_close\"") {
		t.Error("expected save_close action in decoded HTML")
	}
}

func TestBlogSettingsForm_RenderStructure_WithEnvOverride(t *testing.T) {
	infoMessage := "The BLOG_TOPIC environment variable is set, so updates are disabled here."

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	component := NewFormBlogSettings(registry)
	if component == nil {
		t.Fatal("Failed to create component")
	}

	// Set test data on the component
	if comp, ok := component.(*formBlogSettings); ok {
		comp.BlogTopic = "Env controlled topic"
		comp.FormInfoMessage = infoMessage
		comp.IsEnvOverride = true
		comp.ReturnURL = shared.NewLinks().PostManager()
	}

	html := component.Render(context.TODO()).ToHTML()

	if !strings.Contains(html, "Env controlled topic") {
		t.Error("expected Env controlled topic in HTML")
	}
	if !strings.Contains(html, infoMessage) {
		t.Error("expected info message in HTML")
	}
	if !strings.Contains(html, `readonly="true"`) {
		t.Error("expected readonly attribute in HTML")
	}
	if !strings.Contains(html, `disabled="true"`) {
		t.Error("expected disabled attribute in HTML")
	}

	// Buttons should still have data-flux-action attributes but be disabled
	if !strings.Contains(html, `data-flux-action="apply"`) {
		t.Error("expected apply action in HTML")
	}
	if !strings.Contains(html, `data-flux-action="save_close"`) {
		t.Error("expected save_close action in HTML")
	}
	// Verify buttons are disabled
	if !strings.Contains(html, `disabled="true"`) {
		t.Error("expected disabled attribute in HTML")
	}
}

func TestBlogSettingsForm_FlashMessages(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	component := NewFormBlogSettings(registry)
	if component == nil {
		t.Fatal("Failed to create component")
	}

	// Set test data on the component
	if comp, ok := component.(*formBlogSettings); ok {
		comp.BlogTopic = "Contracts"
		comp.FormErrorMessage = "Blog topic is required"
		comp.FormSuccessMessage = "Blog settings saved successfully"
		comp.FormRedirect = "/admin/blog?controller=post-manager"
		comp.FormRedirectDelaySeconds = 5
		comp.ReturnURL = shared.NewLinks().PostManager()
	}

	html := component.Render(context.TODO()).ToHTML()

	if !strings.Contains(html, "Blog topic is required") {
		t.Error("expected error message in HTML")
	}
	if !strings.Contains(html, "Blog settings saved successfully") {
		t.Error("expected success message in HTML")
	}
	if !strings.Contains(html, "/admin/blog?controller=post-manager") {
		t.Error("expected redirect URL in HTML")
	}
}

func htmlDecode(s string) string {
	decoded := html.UnescapeString(s)
	if decoded == s {
		return decoded
	}
	return html.UnescapeString(decoded)
}
