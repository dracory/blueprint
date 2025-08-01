package stats

import (
	"project/internal/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	statsHome := rtr.NewRoute().
		SetName("Admin > Visitor Analytics > Home").
		SetPath(links.ADMIN_STATS).
		SetHandler(StatsController().Handler)

	statsCatchAll := rtr.NewRoute().
		SetName("Admin > Visitor Analytics > Catchall").
		SetPath(links.ADMIN_STATS + links.CATCHALL).
		SetHandler(StatsController().Handler)

	return []rtr.RouteInterface{
		statsHome,
		statsCatchAll,
	}
}
