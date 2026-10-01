// Package layout provides the admin section layout. It wraps the shared
// page scaffold from internal/layouts with an admin-specific navbar and
// chrome, so admin code stays co-located with the controllers that use it.
package layout

import (
	"net/http"
	"project/internal/app"
	"project/internal/layouts"

	baselayouts "github.com/dracory/base/layouts"
)

// New creates an admin page using the shared page scaffold with the
// admin navbar.
func New(app app.AppInterface, r *http.Request, options baselayouts.Options) baselayouts.LayoutInterface {
	return layouts.NewPageLayoutWithConfig(app, r, options, layouts.PageLayoutConfig{
		NavbarFn: navbar,
		IsAdmin:  true,
	})
}
