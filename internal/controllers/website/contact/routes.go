package contact

import (
	"net/http"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

// Routes returns the GET and POST routes for the contact page
func Routes(
	registry registry.RegistryInterface,
) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Website > Contact Controller").
			SetPath(links.CONTACT).
			SetMethod(http.MethodGet).
			SetHTMLHandler(NewContactController(registry).AnyIndex),
	}
}
