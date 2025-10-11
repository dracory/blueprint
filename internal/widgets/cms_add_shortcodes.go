package widgets

import (
	"project/internal/types"

	"github.com/dracory/cmsstore"
)

// CmsAddShortcodes adds the shortcodes to the CMS store
//
// Business Logic:
//   - Check if the CMS store is used
//   - Check if the CMS store is nil
//   - Add the shortcodes to the CMS store
//   - Loaded in the main.go file
//
// Parameters:
//   - None
//
// Returns:
//   - None
func CmsAddShortcodes(app types.AppInterface) {
	if !app.GetConfig().GetCmsStoreUsed() {
		return
	}

	if app.GetCmsStore() == nil {
		return
	}

	shortcodes := []cmsstore.ShortcodeInterface{}

	list := WidgetRegistry(app)

	for _, widget := range list {
		shortcodes = append(shortcodes, widget)
	}

	app.GetCmsStore().AddShortcodes(shortcodes)
}
