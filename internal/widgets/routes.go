package widgets

import (
	"project/internal/types"

	"github.com/dracory/rtr"
)

const PATH_COMMENTABLE = "/widgets/commentable"

func Routes(app types.RegistryInterface) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		// rtr.NewRoute().
		// 	SetName("Website > Commentable Widget").
		// 	SetPath(PATH_COMMENTABLE).
		// 	SetHTMLHandler(NewCommentableWidget().Handler),
	}
}
