package cms

import (
	"project/app/links"
	"project/app/middlewares"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	// Create routes
	routes := []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Website > Widget Controller > Handler").
			SetPath(links.WIDGET).
			SetHTMLHandler(NewWidgetController().Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Home Page").
			SetPath(links.HOME).
			SetHTMLHandler(NewCmsController().Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Catch All Pages").
			SetPath(links.CATCHALL).
			SetHTMLHandler(NewCmsController().Handler),
	}

	// Apply stats middleware to specific routes
	for i, route := range routes {
		if route.GetName() == "Website > Cms > Home Page" || route.GetName() == "Website > Cms > Catch All Pages" {
			routes[i] = route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{middlewares.NewStatsMiddleware()})
		}
	}

	return routes
}
