package widgets

import (
	"project/internal/app"
	"project/internal/links"

	basewidgets "github.com/dracory/base/widgets"
)

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
	return basewidgets.WidgetRegistry(basewidgets.RegistryOptions{
		IsEnvironment: func(environment string) bool {
			config := app.GetConfig()
			return config != nil && config.GetAppEnv() == environment
		},
		WidgetURL: links.Website().Widget,
	})
}
