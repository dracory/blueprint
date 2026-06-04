package widgets

import (
	"project/internal/app"

	"github.com/dracory/rtr"
)

const PATH_COMMENTABLE = "/widgets/commentable"

func Routes(app app.AppInterface) []rtr.RouteInterface {
	return []rtr.RouteInterface{
		// rtr.NewRoute().
		// 	SetName("Website > Commentable Widget").
		// 	SetPath(PATH_COMMENTABLE).
		// 	SetHTMLHandler(NewCommentableWidget().Handler),
	}
}
