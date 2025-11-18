package log_manager

import (
	"net/http"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestLogManagerController_RequiresLogStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		// NOTE: Deliberately no LogStore to trigger the error path
	)

	_, resp, err := test.CallStringEndpoint(http.MethodGet, NewLogManagerController(app).Handler, test.NewRequestOptions{})
	assert.NoError(t, err)
	// When log store is missing, we expect a redirect via flash error helper
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
}

func TestLogManagerController_RendersPlaceholders(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithLogStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err)

	html, resp, err := test.CallStringEndpoint(http.MethodGet, NewLogManagerController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Page should render Liveflux placeholders for both filter and table components.
	assert.Contains(t, html, "admin_log_manager_filter")
	assert.Contains(t, html, "admin_log_manager_table")
}

// Additional scenarios (e.g. specific filter combinations) can be covered by
// component-level tests if needed; for the controller we only ensure that the
// page renders and embeds the expected Liveflux placeholders.
