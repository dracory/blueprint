package admin

import (
	adminBlog "project/app/controllers/admin/blog"
	adminCms "project/app/controllers/admin/cms"
	adminFiles "project/app/controllers/admin/files"
	adminMedia "project/app/controllers/admin/media"
	adminShop "project/app/controllers/admin/shop"
	"project/app/controllers/admin/stats"
	adminTasks "project/app/controllers/admin/tasks"
	adminUsers "project/app/controllers/admin/users"
	"project/internal/links"
	"project/internal/middlewares"

	"github.com/dracory/rtr"
)

// Routes these are the routes for the administrator
func Routes() []rtr.RouteInterface {
	home := rtr.NewRoute().
		SetName("Admin > Home").
		SetPath(links.ADMIN_HOME).
		SetHTMLHandler(NewHomeController().Handler)

	homeCatchAll := rtr.NewRoute().
		SetName("Admin > Catch All").
		SetPath(links.ADMIN_HOME + links.CATCHALL).
		SetHTMLHandler(NewHomeController().Handler)

	adminRoutes := []rtr.RouteInterface{}
	adminRoutes = append(adminRoutes, adminBlog.Routes()...)
	adminRoutes = append(adminRoutes, adminCms.Routes()...)
	adminRoutes = append(adminRoutes, adminFiles.Routes()...)
	adminRoutes = append(adminRoutes, adminMedia.Routes()...)
	adminRoutes = append(adminRoutes, adminShop.ShopRoutes()...)
	adminRoutes = append(adminRoutes, stats.Routes()...)
	adminRoutes = append(adminRoutes, adminTasks.TaskRoutes()...)
	adminRoutes = append(adminRoutes, adminUsers.UserRoutes()...)
	// adminRoutes = append(adminRoutes, []rtr.RouteInterface{subscriptionPlans}...)
	adminRoutes = append(adminRoutes, []rtr.RouteInterface{home, homeCatchAll}...)

	// Apply middlewares to all admin routes
	for _, route := range adminRoutes {
		route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{middlewares.NewAdminMiddleware()})
	}

	return adminRoutes
}
