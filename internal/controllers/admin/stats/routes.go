package stats

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
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
	}
}
