package detailscomponent

import (
	"testing"

	"project/internal/testutils"
)

// TestRenderNotNil verifies Render returns non-nil result
func TestRenderNotNil(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := Render(app, "123")

	if result == nil {
		t.Error("Render() should return non-nil result")
	}
}

// TestRenderWithEmptyCategoryID verifies Render handles empty category ID
func TestRenderWithEmptyCategoryID(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := Render(app, "")

	if result == nil {
		t.Error("Render() should return non-nil result even with empty category ID")
	}
}
