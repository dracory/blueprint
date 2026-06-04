package widgets

import "project/internal/app"

// WidgetRegistry returns a list of all widgets
//
// Register all the new widgets here so that they can be used in the CMS
//
// Parameters:
//   - app: app interface for accessing services
//
// Returns:
//   - []Widget - A list of all widgets
func WidgetRegistry(app app.AppInterface) []Widget {
	return []Widget{
		NewAuthenticatedWidget(app),
		// NewContactFormWidget(app),
		// NewTermsOfUseWidget(app),
		NewUnauthenticatedWidget(app),
		NewVisibleWidget(app),
		NewPrintWidget(app),
		NewBlockeditorWidget(app),
	}
}
