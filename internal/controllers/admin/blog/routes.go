package admin

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {

	postCreate := rtr.NewRoute().
		SetName("Admin > Blog > Post Create").
		SetPath(links.ADMIN_BLOG_POST_CREATE).
		SetHTMLHandler(NewPostCreateController(app).Handler)

	postDelete := rtr.NewRoute().
		SetName("Admin > Blog > Post Delete").
		SetPath(links.ADMIN_BLOG_POST_DELETE).
		SetHTMLHandler(NewPostDeleteController(app).Handler)

	postManager := rtr.NewRoute().
		SetName("Admin > Blog > Post Manager").
		SetPath(links.ADMIN_BLOG_POST_MANAGER).
		SetHTMLHandler(NewManagerController(app).Handler)

	postUpdate := rtr.NewRoute().
		SetName("Admin > Blog > Post Update").
		SetPath(links.ADMIN_BLOG_POST_UPDATE).
		SetHTMLHandler(NewPostUpdateController(app).Handler)

	blogHome := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHTMLHandler(NewManagerController(app).Handler)

	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Blog > Catch All").
		SetPath(links.ADMIN_BLOG + links.CATCHALL).
		SetHTMLHandler(NewManagerController(app).Handler)

	return []rtr.RouteInterface{
		postCreate,
		postDelete,
		postManager,
		postUpdate,
		blogHome,
		blogCatchAll,
	}
}
