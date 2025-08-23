package admin

import (
	adminBlog "project/internal/controllers/admin/blog"
	adminCms "project/internal/controllers/admin/cms"
	adminFiles "project/internal/controllers/admin/files"
	adminMedia "project/internal/controllers/admin/media"
	adminShop "project/internal/controllers/admin/shop"
	"project/internal/controllers/admin/stats"
	adminTasks "project/internal/controllers/admin/tasks"
	adminUsers "project/internal/controllers/admin/users"
	"project/internal/links"
	"project/internal/middlewares"
	"project/internal/types"

	"github.com/dracory/rtr"
)

// Routes these are the routes for the administrator
func Routes(app types.AppInterface) []rtr.RouteInterface {
	home := rtr.NewRoute().
		SetName("Admin > Home").
		SetPath(links.ADMIN_HOME).
		SetHTMLHandler(NewHomeController(app).Handler)

	homeCatchAll := rtr.NewRoute().
		SetName("Admin > Catch All").
		SetPath(links.ADMIN_HOME + links.CATCHALL).
		SetHTMLHandler(NewHomeController(app).Handler)

	adminRoutes := []rtr.RouteInterface{}
	adminRoutes = append(adminRoutes, adminBlog.Routes(app)...)
	adminRoutes = append(adminRoutes, adminCms.Routes(app)...)
	adminRoutes = append(adminRoutes, adminFiles.Routes(app)...)
	adminRoutes = append(adminRoutes, adminMedia.Routes(app)...)
	adminRoutes = append(adminRoutes, adminShop.ShopRoutes(app)...)
	adminRoutes = append(adminRoutes, stats.Routes(app)...)
	adminRoutes = append(adminRoutes, adminTasks.TaskRoutes(app)...)
	adminRoutes = append(adminRoutes, adminUsers.UserRoutes(app)...)
	// adminRoutes = append(adminRoutes, []rtr.RouteInterface{subscriptionPlans}...)
	adminRoutes = append(adminRoutes, []rtr.RouteInterface{home, homeCatchAll}...)

	// Apply middlewares to all admin routes
	for _, route := range adminRoutes {
		route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{middlewares.NewAdminMiddleware(app)})
	}

	return adminRoutes
}
