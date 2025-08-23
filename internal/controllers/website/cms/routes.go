package cms

import (
	"project/internal/links"
	"project/internal/middlewares"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Website > Widget Controller > Handler").
			SetPath(links.WIDGET).
			SetHTMLHandler(NewWidgetController(app).Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Home Page").
			AddBeforeMiddlewares([]rtr.MiddlewareInterface{
				rtr.NewMiddleware().
					SetName("stats").
					SetHandler(middlewares.NewStatsMiddleware(app).GetHandler()),
			}).
			SetPath(links.HOME).
			SetHTMLHandler(NewCmsController(app).Handler),

		rtr.NewRoute().
			SetName("Website > Cms > Catch All Pages").
			AddBeforeMiddlewares([]rtr.MiddlewareInterface{
				rtr.NewMiddleware().
					SetName("stats").
					SetHandler(middlewares.NewStatsMiddleware(app).GetHandler()),
			}).
			SetPath(links.CATCHALL).
			SetHTMLHandler(NewCmsController(app).Handler),
	}
}
