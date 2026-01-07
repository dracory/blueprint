package cms

import (
	"project/internal/links"
	"project/internal/middlewares"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Website > Widget Controller > Handler").
			SetPath(links.WIDGET).
			SetHTMLHandler(NewWidgetController(registry).Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Home Page").
			AddBeforeMiddlewares([]rtr.MiddlewareInterface{
				rtr.NewMiddleware().
					SetName("stats").
					SetHandler(middlewares.NewStatsMiddleware(registry).GetHandler()),
			}).
			SetPath(links.HOME).
			SetHTMLHandler(NewCmsController(registry).Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Catch All Pages").
			AddBeforeMiddlewares([]rtr.MiddlewareInterface{
				rtr.NewMiddleware().
					SetName("stats").
					SetHandler(middlewares.NewStatsMiddleware(registry).GetHandler()),
			}).
			SetPath(links.CATCHALL).
			SetHTMLHandler(NewCmsController(registry).Handler),
	}
}
