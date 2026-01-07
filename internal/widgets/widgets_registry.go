package widgets

import "project/internal/registry"

// WidgetRegistry returns a list of all widgets
//
// Register all the new widgets here so that they can be used in the CMS
//
// Parameters:
//   - None
//
// Returns:
//   - []Widget - A list of all widgets
func WidgetRegistry(registry registry.RegistryInterface) []Widget {
	return []Widget{
		NewAuthenticatedWidget(registry),
		// NewContactFormWidget(registry),
		NewUnauthenticatedWidget(registry),
		NewVisibleWidget(registry),
	}
}
