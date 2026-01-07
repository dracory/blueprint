package admin

import (
	adminBlog "project/internal/controllers/admin/blog"
	adminCms "project/internal/controllers/admin/cms"
	adminFiles "project/internal/controllers/admin/files"
	adminLogs "project/internal/controllers/admin/logs"
	adminMedia "project/internal/controllers/admin/media"
	adminShop "project/internal/controllers/admin/shop"
	adminStats "project/internal/controllers/admin/stats"
	adminTasks "project/internal/controllers/admin/tasks"
	adminUsers "project/internal/controllers/admin/users"
	"project/internal/links"
	"project/internal/middlewares"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

// Routes these are the routes for the administrator
func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	home := rtr.NewRoute().
		SetName("Admin > Home").
		SetPath(links.ADMIN_HOME).
		SetHTMLHandler(NewHomeController(registry).Handler)

	homeCatchAll := rtr.NewRoute().
		SetName("Admin > Catch All").
		SetPath(links.ADMIN_HOME + links.CATCHALL).
		SetHTMLHandler(NewHomeController(registry).Handler)

	adminRoutes := []rtr.RouteInterface{}
	adminRoutes = append(adminRoutes, adminBlog.Routes(registry)...)
	adminRoutes = append(adminRoutes, adminCms.Routes(registry)...)
	adminRoutes = append(adminRoutes, adminFiles.Routes(registry)...)
	adminRoutes = append(adminRoutes, adminLogs.Routes(registry)...)
	adminRoutes = append(adminRoutes, adminMedia.Routes(registry)...)
	adminRoutes = append(adminRoutes, adminShop.ShopRoutes(registry)...)
	adminRoutes = append(adminRoutes, adminStats.Routes(registry)...)
	adminRoutes = append(adminRoutes, adminTasks.TaskRoutes(registry)...)
	adminRoutes = append(adminRoutes, adminUsers.UserRoutes(registry)...)
	// adminRoutes = append(adminRoutes, []rtr.RouteInterface{subscriptionPlans}...)
	adminRoutes = append(adminRoutes, []rtr.RouteInterface{home, homeCatchAll}...)

	// Apply middlewares to all admin routes
	for _, route := range adminRoutes {
		route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			middlewares.NewAdminMiddleware(registry),
		})
	}

	return adminRoutes
}
