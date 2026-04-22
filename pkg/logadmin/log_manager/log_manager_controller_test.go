package log_manager

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/test"
)

func TestLogManagerController_RendersVueApp(t *testing.T) {
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
	// Page should render Vue app mount point
	if !strings.Contains(html, "logs-app") {
		t.Error("expected logs-app in HTML")
	}
	// Page should include Vue CDN
	if !strings.Contains(html, "vue.global.js") {
		t.Error("expected Vue CDN script in HTML")
	}
	// Page should include SweetAlert2
	if !strings.Contains(html, "sweetalert2") {
		t.Error("expected SweetAlert2 in HTML")
	}
}

func TestLogManagerController_LoadLogsAction(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithLogStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test load-logs action
	queryParams := url.Values{}
	queryParams.Set("action", "load-logs")

	requestBody := map[string]any{
		"page":       0,
		"per_page":   100,
		"sort_order": "desc",
		"sort_by":    "time",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	body, resp, err := test.CallStringEndpoint(http.MethodPost, NewLogManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: queryParams,
		Body:      string(bodyBytes),
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
	// Should return JSON response
	if !strings.Contains(body, `"status"`) {
		t.Error("expected JSON response with status field")
	}
	if !strings.Contains(body, `"logs"`) {
		t.Error("expected JSON response with logs field")
	}
}
