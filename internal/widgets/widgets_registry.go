package widgets

import "project/internal/types"

// WidgetRegistry returns a list of all widgets
//
// Register all the new widgets here so that they can be used in the CMS
//
// Parameters:
//   - None
//
// Returns:
//   - []Widget - A list of all widgets
func WidgetRegistry(app types.RegistryInterface) []Widget {
	return []Widget{
		NewAuthenticatedWidget(app),
		// NewContactFormWidget(app),
		NewUnauthenticatedWidget(app),
		NewVisibleWidget(app),
	}
}
