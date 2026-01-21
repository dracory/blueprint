package widgets

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/registry"
)

var _ Widget = (*authenticatedWidget)(nil) // verify it extends the interface

// == CONSTUCTOR ==============================================================

// NewAuthenticatedWidget returns a new instance of authenticatedWidget.
//
// Parameters:
// - registry: Registry interface for accessing services
//
// Returns:
// - *authenticatedWidget: a pointer to an authenticatedWidget.
func NewAuthenticatedWidget(registry registry.RegistryInterface) *authenticatedWidget {
	return &authenticatedWidget{registry: registry}
}

// == WIDGET =================================================================

// authenticatedWidget used to render the authenticatedWidget shortcode.
//
// It displays the content of the shortcode if the user is authenticated.
type authenticatedWidget struct {
	registry registry.RegistryInterface
}

// == PUBLIC METHODS =========================================================

// Alias the shortcode alias to be used in the template.
func (w *authenticatedWidget) Alias() string {
	return "x-authenticated"
}

// Description a user-friendly description of the shortcode.
func (w *authenticatedWidget) Description() string {
	return "Renders the content if the user is authenticated"
}

// Render implements the shortcode interface.
func (w *authenticatedWidget) Render(req *http.Request, content string, data map[string]string) string {
	authUser := helpers.GetAuthUser(req)

	if authUser != nil {
		return content
	}

	return ""
}
