package widgets

import (
	"project/internal/registry"

	"github.com/dracory/cmsstore"
)

// CmsAddShortcodes adds the shortcodes to the CMS store
//
// Business Logic:
//   - Check if the CMS store is used
//   - Check if the CMS store is nil
//   - Add the shortcodes to the CMS store
//   - Loaded in the cmd/server/main.go file
//
// Parameters:
//   - None
//
// Returns:
//   - None
func CmsAddShortcodes(registry registry.RegistryInterface) {
	if !registry.GetConfig().GetCmsStoreUsed() {
		return
	}

	if registry.GetCmsStore() == nil {
		return
	}

	shortcodes := []cmsstore.ShortcodeInterface{}

	list := WidgetRegistry(registry)

	for _, widget := range list {
		shortcodes = append(shortcodes, widget)
	}

	registry.GetCmsStore().AddShortcodes(shortcodes)
}
