package contact

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

// Routes returns the GET and POST routes for the contact page
func Routes(
	app types.AppInterface,
) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Website > Contact Controller").
			SetPath(links.CONTACT).
			SetMethod(http.MethodGet).
			SetHTMLHandler(NewContactController(app).AnyIndex),
	}
}
