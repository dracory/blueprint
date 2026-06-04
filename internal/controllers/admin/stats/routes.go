package stats

import (
	"errors"
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/rtr"
)

func Routes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	ctrl := NewStatsController(app)
	statsHome := rtr.NewRoute().
		SetName("Admin > Visitor Analytics > Home").
		SetPath(links.ADMIN_STATS).
		SetHandler(ctrl.Handler)

	statsCatchAll := rtr.NewRoute().
		SetName("Admin > Visitor Analytics > Catchall").
		SetPath(links.ADMIN_STATS + links.CATCHALL).
		SetHandler(ctrl.Handler)

	return []rtr.RouteInterface{
		statsHome,
		statsCatchAll,
	}, nil
}
