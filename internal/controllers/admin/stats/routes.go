package stats

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	ctrl := NewStatsController(registry)
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
