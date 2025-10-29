package admin

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"

	postCreate "project/internal/controllers/admin/blog/post_create"
	postDelete "project/internal/controllers/admin/blog/post_delete"
	postManager "project/internal/controllers/admin/blog/post_manager"
	postUpdate "project/internal/controllers/admin/blog/post_update"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {

	postCreateRoute := rtr.NewRoute().
		SetName("Admin > Blog > Post Create").
		SetPath(links.ADMIN_BLOG_POST_CREATE).
		SetHTMLHandler(postCreate.NewPostCreateController(app).Handler)

	postDeleteRoute := rtr.NewRoute().
		SetName("Admin > Blog > Post Delete").
		SetPath(links.ADMIN_BLOG_POST_DELETE).
		SetHTMLHandler(postDelete.NewPostDeleteController(app).Handler)

	postManagerRoute := rtr.NewRoute().
		SetName("Admin > Blog > Post Manager").
		SetPath(links.ADMIN_BLOG_POST_MANAGER).
		SetHTMLHandler(postManager.NewManagerController(app).Handler)

	postUpdateRoute := rtr.NewRoute().
		SetName("Admin > Blog > Post Update").
		SetPath(links.ADMIN_BLOG_POST_UPDATE).
		SetHTMLHandler(postUpdate.NewPostUpdateController(app).Handler)

	blogHome := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHTMLHandler(postManager.NewManagerController(app).Handler)

	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Blog > Catch All").
		SetPath(links.ADMIN_BLOG + links.CATCHALL).
		SetHTMLHandler(postManager.NewManagerController(app).Handler)

	return []rtr.RouteInterface{
		postCreateRoute,
		postDeleteRoute,
		postManagerRoute,
		postUpdateRoute,
		blogHome,
		blogCatchAll,
	}
}
