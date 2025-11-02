package blog_settings

import (
	"html"
	"strings"
	"testing"

	"project/internal/controllers/admin/blog/shared"

	"github.com/stretchr/testify/assert"
)

func TestBlogSettingsForm_RenderStructure(t *testing.T) {
	options := blogSettingsFormOptions{Data: blogSettingsData{
		blogTopic: "AI insights",
	}}

	html := blogSettingsForm(options).ToHTML()
	decoded := htmlDecode(html)
	t.Log(decoded)

	assert.Contains(t, html, `id="BlogSettingsFormWrapper"`)
	assert.Contains(t, html, `name="blog_topic"`)
	assert.Contains(t, html, "AI insights")
	cancelURL := shared.NewLinks().PostManager()
	assert.Contains(t, html, cancelURL)

	blogSettingsURL := shared.NewLinks().BlogSettings()
	assert.Contains(t, html, `hx-post="`+blogSettingsURL+`"`)
	assert.Contains(t, html, `hx-target="#BlogSettingsFormWrapper"`)
	assert.Contains(t, html, `hx-indicator="#BlogSettingsApplyIndicator"`)
	assert.Contains(t, html, `hx-indicator="#BlogSettingsSaveIndicator"`)
	assert.Contains(t, html, `id="BlogSettingsApplyIndicator"`)
	assert.Contains(t, html, `id="BlogSettingsSaveIndicator"`)

	assert.Contains(t, decoded, "\"action\":\"apply\"")
	assert.Contains(t, decoded, "\"action\":\"save_close\"")
}

func TestBlogSettingsForm_RenderStructure_WithEnvOverride(t *testing.T) {
	infoMessage := "The BLOG_TOPIC environment variable is set, so updates are disabled here."
	options := blogSettingsFormOptions{Data: blogSettingsData{
		blogTopic:       "Env controlled topic",
		formInfoMessage: infoMessage,
		isEnvOverride:   true,
	}}

	html := blogSettingsForm(options).ToHTML()

	assert.Contains(t, html, "Env controlled topic")
	assert.Contains(t, html, infoMessage)
	assert.Contains(t, html, `readonly="true"`)
	assert.Contains(t, html, `disabled="true"`)

	assert.NotContains(t, html, `hx-vals="{\"action\":\"apply\"}"`)
	assert.NotContains(t, html, `hx-vals="{\"action\":\"save_close\"}"`)
	assert.Contains(t, html, `disabled="true"`)
}

func TestBlogSettingsForm_FlashMessages(t *testing.T) {
	options := blogSettingsFormOptions{Data: blogSettingsData{
		blogTopic:                "Contracts",
		formErrorMessage:         "Blog topic is required",
		formSuccessMessage:       "Blog settings saved successfully",
		formRedirect:             "/admin/blog?controller=post-manager",
		formRedirectDelaySeconds: 5,
	}}

	html := blogSettingsForm(options).ToHTML()

	assert.Contains(t, html, "Blog topic is required")
	assert.Contains(t, html, "Blog settings saved successfully")
	assert.Contains(t, html, options.Data.formRedirect)
}

func containsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func htmlDecode(s string) string {
	decoded := html.UnescapeString(s)
	if decoded == s {
		return decoded
	}
	return html.UnescapeString(decoded)
}
