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
func WidgetRegistry(app types.AppInterface) []Widget {
	return []Widget{
		NewAuthenticatedWidget(app),
		NewContactFormWidget(app),
		NewTermsOfUseWidget(app),
		NewUnauthenticatedWidget(app),
		NewVisibleWidget(app),
	}
}
