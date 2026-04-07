package log_manager

import (
	"net/http"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestLogManagerController_RequiresLogStore(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		// NOTE: Deliberately no LogStore to trigger the error path
	)

	_, resp, err := test.CallStringEndpoint(http.MethodGet, NewLogManagerController(registry).Handler, test.NewRequestOptions{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// When log store is missing, we expect a redirect via flash error helper
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, resp.StatusCode)
	}
}

func TestLogManagerController_RendersPlaceholders(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithLogStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	html, resp, err := test.CallStringEndpoint(http.MethodGet, NewLogManagerController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	// Page should render Liveflux placeholders for both filter and table components.
	if !strings.Contains(html, "admin_log_manager_filter") {
		t.Error("expected admin_log_manager_filter in HTML")
	}
	if !strings.Contains(html, "admin_log_manager_table") {
		t.Error("expected admin_log_manager_table in HTML")
	}
}

// Additional scenarios (e.g. specific filter combinations) can be covered by
// component-level tests if needed; for the controller we only ensure that the
// page renders and embeds the expected Liveflux placeholders.
