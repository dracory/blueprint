package blog_settings

import (
	"context"
	"html"
	"testing"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/testutils"

	"github.com/stretchr/testify/assert"
)

func TestBlogSettingsForm_RenderStructure(t *testing.T) {
	// Create a mock app for testing
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	component := NewFormBlogSettings(app)
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

	assert.Contains(t, html, `id="BlogSettingsFormWrapper"`)
	assert.Contains(t, html, `name="blog_topic"`)
	assert.Contains(t, html, "AI insights")
	cancelURL := shared.NewLinks().PostManager()
	assert.Contains(t, html, cancelURL)

	assert.Contains(t, html, `data-flux-action="apply"`)
	assert.Contains(t, html, `data-flux-action="save_close"`)

	// Verify the action attributes are present in the HTML
	assert.Contains(t, decoded, "data-flux-action=\"apply\"")
	assert.Contains(t, decoded, "data-flux-action=\"save_close\"")
}

func TestBlogSettingsForm_RenderStructure_WithEnvOverride(t *testing.T) {
	infoMessage := "The BLOG_TOPIC environment variable is set, so updates are disabled here."

	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	component := NewFormBlogSettings(app)
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

	assert.Contains(t, html, "Env controlled topic")
	assert.Contains(t, html, infoMessage)
	assert.Contains(t, html, `readonly="true"`)
	assert.Contains(t, html, `disabled="true"`)

	// Buttons should still have data-flux-action attributes but be disabled
	assert.Contains(t, html, `data-flux-action="apply"`)
	assert.Contains(t, html, `data-flux-action="save_close"`)
	// Verify buttons are disabled
	assert.Contains(t, html, `disabled="true"`)
}

func TestBlogSettingsForm_FlashMessages(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	component := NewFormBlogSettings(app)
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

	assert.Contains(t, html, "Blog topic is required")
	assert.Contains(t, html, "Blog settings saved successfully")
	assert.Contains(t, html, "/admin/blog?controller=post-manager")
}

func htmlDecode(s string) string {
	decoded := html.UnescapeString(s)
	if decoded == s {
		return decoded
	}
	return html.UnescapeString(decoded)
}
