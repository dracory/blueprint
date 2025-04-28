package stats

import (
	"project/app/links"

	"github.com/gouniverse/router"
)

func Routes() []router.RouteInterface {
	return []router.RouteInterface{
		&router.Route{
			Name:    "Admin > Visitor Analytics > Home",
			Path:    links.ADMIN_STATS,
			Handler: StatsController().Handler,
		},
		&router.Route{
			Name:    "Admin > Visitor Analytics > Catchall",
			Path:    links.ADMIN_STATS + links.CATCHALL,
			Handler: StatsController().Handler,
		},
	}
}
