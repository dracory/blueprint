package liveflux

import (
	"net/http"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

// Routes registers the /liveflux endpoint for the Liveflux handler.
func Routes(app types.AppInterface) []rtr.RouteInterface {
	ctrl := NewController(app)

	livefluxPost := rtr.NewRoute().
		SetName("Liveflux > Handler").
		SetPath(links.LIVEFLUX).
		SetMethod(http.MethodPost).
		SetHTMLHandler(ctrl.Handler)

	livefluxGet := rtr.NewRoute().
		SetName("Liveflux > Handler").
		SetPath(links.LIVEFLUX).
		SetMethod(http.MethodGet).
		SetHTMLHandler(ctrl.Handler)

	return []rtr.RouteInterface{livefluxPost, livefluxGet}
}
