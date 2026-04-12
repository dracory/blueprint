package widgets

import (
	"testing"

	"project/internal/testutils"
)

// TestCmsAddShortcodes tests the CmsAddShortcodes function
func TestCmsAddShortcodes(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup()

	// Should not panic when CMS store is nil or not used
	CmsAddShortcodes(registry)
}

// TestCmsAddShortcodes_MultipleCalls tests calling the function multiple times
func TestCmsAddShortcodes_MultipleCalls(t *testing.T) {
	t.Parallel()

	registry := testutils.Setup()

	// Should be safe to call multiple times without panic
	CmsAddShortcodes(registry)
	CmsAddShortcodes(registry)
	CmsAddShortcodes(registry)
}
