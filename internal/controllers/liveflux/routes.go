package liveflux

import (
	"net/http"

	"github.com/dracory/rtr"
)

// Routes registers the /liveflux endpoint for the Liveflux handler.
func Routes() []rtr.RouteInterface {
	ctrl := NewController()

	r := rtr.NewRoute().
		SetName("Liveflux > Handler").
		SetPath("/liveflux").
		SetMethod(http.MethodPost).
		SetHTMLHandler(ctrl.Handler)

	return []rtr.RouteInterface{r}
}
