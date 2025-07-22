package stats

import (
	"project/app/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	return []rtr.RouteInterface{
		rtr.NewRoute().
			SetName("Admin > Visitor Analytics > Home").
			SetPath(links.ADMIN_STATS).
			SetHandler(StatsController().Handler),
		rtr.NewRoute().
			SetName("Admin > Visitor Analytics > Catchall").
			SetPath(links.ADMIN_STATS + links.CATCHALL).
			SetHandler(StatsController().Handler),
	}
}
