package liveflux

import (
	"net/http"
	"project/internal/types"

	"github.com/dracory/rtr"
)

// Routes registers the /liveflux endpoint for the Liveflux handler.
func Routes(app types.AppInterface) []rtr.RouteInterface {
	ctrl := NewController(app)

	r := rtr.NewRoute().
		SetName("Liveflux > Handler").
		SetPath("/liveflux").
		SetMethod(http.MethodPost).
		SetHTMLHandler(ctrl.Handler)

	return []rtr.RouteInterface{r}
}
