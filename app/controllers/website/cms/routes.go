package cms

import (
	"project/app/links"
	"project/app/middlewares"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Website > Widget Controller > Handler").
			SetPath(links.WIDGET).
			SetHTMLHandler(NewWidgetController().Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Home Page").
			AddBeforeMiddlewares([]rtr.MiddlewareInterface{
				rtr.NewMiddleware().
					SetName("stats").
					SetHandler(middlewares.NewStatsMiddleware().GetHandler()),
			}).
			SetPath(links.HOME).
			SetHTMLHandler(NewCmsController().Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Catch All Pages").
			AddBeforeMiddlewares([]rtr.MiddlewareInterface{
				rtr.NewMiddleware().
					SetName("stats").
					SetHandler(middlewares.NewStatsMiddleware().GetHandler()),
			}).
			SetPath(links.CATCHALL).
			SetHTMLHandler(NewCmsController().Handler),
	}
}
