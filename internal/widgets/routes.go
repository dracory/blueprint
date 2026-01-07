package widgets

import (
	"project/internal/registry"

	"github.com/dracory/rtr"
)

const PATH_COMMENTABLE = "/widgets/commentable"

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		// rtr.NewRoute().
		// 	SetName("Website > Commentable Widget").
		// 	SetPath(PATH_COMMENTABLE).
		// 	SetHTMLHandler(NewCommentableWidget().Handler),
	}
}
