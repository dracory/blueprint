package admin

import (
	"net/http"
	"project/internal/app"
	adminBlog "project/internal/controllers/admin/blog"
	adminCms "project/internal/controllers/admin/cms"
	adminFiles "project/internal/controllers/admin/files"
	adminMedia "project/internal/controllers/admin/media"
	adminShop "project/internal/controllers/admin/shop"
	adminStats "project/internal/controllers/admin/stats"
	adminTasks "project/internal/controllers/admin/tasks"
	adminUsers "project/internal/controllers/admin/users"
	"project/internal/links"
	"project/internal/middlewares"
	"project/pkg/logadmin"

	"github.com/dracory/rtr"
)

// Routes these are the routes for the administrator
func Routes(app app.AppInterface) []rtr.RouteInterface {
	home := rtr.NewRoute().
		SetName("Admin > Home").
		SetPath(links.ADMIN_HOME).
		SetHTMLHandler(NewHomeController(app).Handler)

	homeCatchAll := rtr.NewRoute().
		SetName("Admin > Catch All").
		SetPath(links.ADMIN_HOME + links.CATCHALL).
		SetHTMLHandler(NewHomeController(app).Handler)

	adminRoutes := []rtr.RouteInterface{}

	blogController := adminBlog.NewBlogAdminController(app)
	blog := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			blogController.Handler(w, r)
			return ""
		})
	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Blog > Catchall").
		SetPath(links.ADMIN_BLOG + links.CATCHALL).
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			blogController.Handler(w, r)
			return ""
		})
	adminRoutes = append(adminRoutes, blog, blogCatchAll)

	cmsRoutes, err := adminCms.Routes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, cmsRoutes...)
	}

	fileRoutes, err := adminFiles.Routes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, fileRoutes...)
	}

	logRoutes, err := logadmin.Routes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, logRoutes...)
	}

	mediaRoutes, err := adminMedia.Routes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, mediaRoutes...)
	}

	shopController := adminShop.NewShopAdminController(app)
	shop := rtr.NewRoute().
		SetName("Admin > Shop").
		SetPath(links.ADMIN_SHOP).
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			shopController.Handler(w, r)
			return ""
		})
	shopCatchAll := rtr.NewRoute().
		SetName("Admin > Shop > Catchall").
		SetPath(links.ADMIN_SHOP + links.CATCHALL).
		SetHTMLHandler(func(w http.ResponseWriter, r *http.Request) string {
			shopController.Handler(w, r)
			return ""
		})
	adminRoutes = append(adminRoutes, shop, shopCatchAll)

	statsRoutes, err := adminStats.Routes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, statsRoutes...)
	}

	taskRoutes, err := adminTasks.TaskRoutes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, taskRoutes...)
	}

	userRoutes, err := adminUsers.Routes(app)
	if err == nil {
		adminRoutes = append(adminRoutes, userRoutes...)
	}

	// adminRoutes = append(adminRoutes, []rtr.RouteInterface{subscriptionPlans}...)
	adminRoutes = append(adminRoutes, []rtr.RouteInterface{home, homeCatchAll}...)

	// Apply middlewares to all admin routes
	for _, route := range adminRoutes {
		route.AddBeforeMiddlewares([]rtr.MiddlewareInterface{
			middlewares.NewAdminMiddleware(app),
			middlewares.NewEmailAllowlistMiddleware(app),
		})
	}

	return adminRoutes
}
