package swagger

import (
	"net/http"

	"github.com/dracory/rtr"
)

// Routes returns the routes for Swagger UI and YAML
func Routes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Swagger UI").
			SetPath("/swagger").
			SetHandler(SwaggerUIController).
			SetMethod(http.MethodGet),

		rtr.NewRoute().
			SetName("Swagger YAML").
			SetPath("/docs/swagger.yaml").
			SetHandler(SwaggerYAMLController).
			SetMethod(http.MethodGet),
	}
}
