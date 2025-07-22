package admin

import (
	"project/app/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {

	postCreate := rtr.NewRoute().
		SetName("Admin > Blog > Post Create").
		SetPath(links.ADMIN_BLOG_POST_CREATE).
		SetHTMLHandler(NewPostCreateController().Handler)

	postDelete := rtr.NewRoute().
		SetName("Admin > Blog > Post Delete").
		SetPath(links.ADMIN_BLOG_POST_DELETE).
		SetHTMLHandler(NewPostDeleteController().Handler)

	postManager := rtr.NewRoute().
		SetName("Admin > Blog > Post Manager").
		SetPath(links.ADMIN_BLOG_POST_MANAGER).
		SetHTMLHandler(NewManagerController().Handler)

	postUpdate := rtr.NewRoute().
		SetName("Admin > Blog > Post Update").
		SetPath(links.ADMIN_BLOG_POST_UPDATE).
		SetHTMLHandler(NewPostUpdateController().Handler)

	blogHome := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHTMLHandler(NewManagerController().Handler)

	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Blog > Catch All").
		SetPath(links.ADMIN_BLOG + links.CATCHALL).
		SetHTMLHandler(NewManagerController().Handler)

	return []rtr.RouteInterface{
		postCreate,
		postDelete,
		postManager,
		postUpdate,
		blogHome,
		blogCatchAll,
	}
}
