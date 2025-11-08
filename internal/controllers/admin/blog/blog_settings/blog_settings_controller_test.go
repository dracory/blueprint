package blog_settings

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"project/internal/config"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/testutils"

	"github.com/dracory/cdn"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestBlogSettingsController_Handler_RendersAssets(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err)

	// Seed existing value to ensure store is operational
	assert.NoError(t, app.GetSettingStore().Set(context.Background(), SettingKeyBlogTopic, "Seeded Topic"))

	html, resp, err := test.CallStringEndpoint(http.MethodGet, NewBlogSettingsController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, html, "Blog Settings")
	assert.Contains(t, html, shared.NewLinks().BlogSettings())
	assert.Contains(t, html, cdn.Htmx_2_0_0())
	assert.Contains(t, html, cdn.Sweetalert2_11())
}

func TestBlogSettingsController_Handler_WithEnvOverride(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err)

	const envValue = "Env Topic"
	os.Setenv("BLOG_TOPIC", envValue)
	t.Cleanup(func() { os.Unsetenv("BLOG_TOPIC") })

	// GET should render env override messaging and disable inputs
	getHTML, resp, err := test.CallStringEndpoint(http.MethodGet, NewBlogSettingsController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, getHTML, envValue)
	assert.Contains(t, getHTML, "updates are disabled here")
	assert.Contains(t, getHTML, "readonly=\"true\"")
	assert.Contains(t, getHTML, "disabled=\"true\"")

	// POST should not mutate the store and should show error message
	postHTML, postResp, err := test.CallStringEndpoint(http.MethodPost, NewBlogSettingsController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
		PostValues: url.Values{
			"blog_topic": {"Attempted Update"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	assert.Contains(t, postHTML, "Blog topic is managed via environment and cannot be changed here.")
	assert.Contains(t, postHTML, envValue)
}
