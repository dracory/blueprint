package shared

import (
	"net/http"
	"project/internal/app"
	"project/internal/controllers/shared/flash"
	"project/internal/controllers/shared/page_not_found"
	"project/internal/controllers/shared/theme"
	"project/internal/links"
	"project/internal/resources"

	"github.com/dracory/base/fileserve"

	"github.com/dracory/rtr"
	"github.com/dracory/uncdn"
)

func Routes(app app.AppInterface) []rtr.RouteInterface {
	cdnRoute := rtr.NewRoute().
		SetName("Shared > CDN Controller").
		SetPath("/cdn/:name").
		SetMethod(http.MethodGet).
		SetHandler(uncdn.CdnHandler())

	files := rtr.NewRoute().
		SetName("Shared > Files Controller").
		SetPath(links.FILES).
		SetMethod(http.MethodGet).
		SetHTMLHandler(fileserve.FileHandler(app.GetSqlFileStorage()))

	flash := rtr.NewRoute().
		SetName("Shared > Flash Controller").
		SetPath(links.FLASH).
		SetHTMLHandler(flash.NewFlashController(app).Handler)

	media := rtr.NewRoute().
		SetName("Shared > Media Controller").
		SetPath(links.MEDIA).
		SetMethod(http.MethodGet).
		SetHTMLHandler(fileserve.FileHandler(app.GetSqlFileStorage()))

	resourcesRoute := rtr.NewRoute().
		SetName("Shared > Resources Controller").
		SetPath(links.RESOURCES).
		SetHTMLHandler(fileserve.ResourceHandler(resources.Resource, page_not_found.PageNotFoundController(nil).Handler))

	themeRoute := rtr.NewRoute().
		SetName("Shared > Theme Controller").
		SetPath(links.THEME).
		SetHandler(theme.ThemeController)

	thumbRoute := rtr.NewRoute().
		SetName("Shared > Thumb Controller").
		SetPath("/th/:extension/:size/:quality/:path").
		SetHTMLHandler(fileserve.ThumbHandler(fileserve.ThumbOptions{
			Logger:         app.GetLogger(),
			Cache:          app.GetFileCache(),
			Storage:        app.GetSqlFileStorage(),
			ResourceReader: resources.ToBytes,
		}))

	thumbRoutePathCarchAll := rtr.NewRoute().
		SetName("Shared > Thumb Controller").
		SetPath("/th/:extension/:size/:quality/:path...").
		SetHTMLHandler(fileserve.ThumbHandler(fileserve.ThumbOptions{
			Logger:         app.GetLogger(),
			Cache:          app.GetFileCache(),
			Storage:        app.GetSqlFileStorage(),
			ResourceReader: resources.ToBytes,
		}))

	thumbCatchAll := rtr.NewRoute().
		SetName("Shared > Thumb Controller").
		SetPath("/th/*").
		SetHTMLHandler(fileserve.ThumbHandler(fileserve.ThumbOptions{
			Logger:         app.GetLogger(),
			Cache:          app.GetFileCache(),
			Storage:        app.GetSqlFileStorage(),
			ResourceReader: resources.ToBytes,
		}))

	pageNotFound := rtr.NewRoute().
		SetName("Shared > Page Not Found Controller").
		SetPath("/404").
		SetHTMLHandler(page_not_found.PageNotFoundController(app).Handler)

	return []rtr.RouteInterface{
		cdnRoute,
		files,
		flash,
		media,
		resourcesRoute,
		themeRoute,
		thumbRoute,
		thumbRoutePathCarchAll,
		thumbCatchAll,
		pageNotFound,
	}
}
