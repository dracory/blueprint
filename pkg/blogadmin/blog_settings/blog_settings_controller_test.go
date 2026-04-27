package blog_settings

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/links"
	"project/internal/testutils"
	"project/pkg/blogadmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/test"
)

func TestBlogSettingsController_Handler_RendersAssets(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("SeedUser returned error: %v", err)
	}

	// Seed existing value to ensure store is operational
	if err := registry.GetSettingStore().Set(context.Background(), SettingKeyBlogTopic, "Seeded Topic"); err != nil {
		t.Errorf("Set returned error: %v", err)
	}

	html, resp, err := test.CallStringEndpoint(http.MethodGet, NewBlogSettingsController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if !strings.Contains(html, "Blog Settings") {
		t.Error("expected Blog Settings in HTML")
	}
	if !strings.Contains(html, shared.NewLinks("/admin/blog").BlogSettings()) {
		t.Error("expected BlogSettings link in HTML")
	}
	if !strings.Contains(html, links.LIVEFLUX) {
		t.Error("expected LIVEFLUX in HTML")
	}
	if !strings.Contains(html, cdn.Sweetalert2_11()) {
		t.Error("expected Sweetalert2 CDN in HTML")
	}
}

func TestBlogSettingsController_Handler_WithEnvOverride(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithSettingStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("SeedUser returned error: %v", err)
	}

	const envValue = "Env Topic"
	os.Setenv("BLOG_TOPIC", envValue)
	t.Cleanup(func() { os.Unsetenv("BLOG_TOPIC") })

	// Seed existing value to ensure store is operational
	if err := registry.GetSettingStore().Set(context.Background(), SettingKeyBlogTopic, "Seeded Topic"); err != nil {
		t.Errorf("Set returned error: %v", err)
	}

	// GET should render env override messaging and disable inputs
	getHTML, resp, err := test.CallStringEndpoint(http.MethodGet, NewBlogSettingsController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if !strings.Contains(getHTML, envValue) {
		t.Error("expected env value in HTML")
	}
	if !strings.Contains(getHTML, "updates are disabled here") {
		t.Error("expected disabled message in HTML")
	}
	if !strings.Contains(getHTML, "readonly=\"true\"") {
		t.Error("expected readonly attribute in HTML")
	}
	if !strings.Contains(getHTML, "disabled=\"true\"") {
		t.Error("expected disabled attribute in HTML")
	}

	// POST should not mutate the store and should show error message
	postHTML, postResp, err := test.CallStringEndpoint(http.MethodPost, NewBlogSettingsController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
		PostValues: url.Values{
			"blog_topic": {"Attempted Update"},
		},
	})

	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
	if postResp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, postResp.StatusCode)
	}
	if !strings.Contains(postHTML, "Blog topic is managed via environment and cannot be changed here.") {
		t.Error("expected error message in HTML")
	}
	if !strings.Contains(postHTML, envValue) {
		t.Error("expected env value in HTML")
	}
}
