package shared

import (
	"net/http"
	"project/internal/controllers/shared/file"
	"project/internal/controllers/shared/flash"
	"project/internal/controllers/shared/media"
	"project/internal/controllers/shared/resource"
	"project/internal/controllers/shared/thumb"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
	"github.com/gouniverse/dashboard"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	adsTxt := rtr.NewRoute().
		SetName("Shared > ads.txt").
		SetPath("/ads.txt").
		SetStringHandler(func(w http.ResponseWriter, r *http.Request) string {
			return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
		})

	files := rtr.NewRoute().
		SetName("Shared > Files Controller").
		SetPath(links.FILES).
		SetMethod(http.MethodGet).
		SetHTMLHandler(file.NewFileController(app.GetSqlFileStorage()).Handler)

	flash := rtr.NewRoute().
		SetName("Shared > Flash Controller").
		SetPath(links.FLASH).
		SetHTMLHandler(flash.NewFlashController(app).Handler)

	media := rtr.NewRoute().
		SetName("Shared > Media Controller").
		SetPath(links.MEDIA).
		SetMethod(http.MethodGet).
		SetHTMLHandler(media.NewMediaController(app.GetSqlFileStorage()).Handler)

	resources := rtr.NewRoute().
		SetName("Shared > Resources Controller").
		SetPath(links.RESOURCES).
		SetHTMLHandler(resource.NewResourceController().Handler)

	theme := rtr.NewRoute().
		SetName("Shared > Theme Controller").
		SetPath(links.THEME).
		SetHandler(dashboard.ThemeHandler)

	thumb := rtr.NewRoute().
		SetName("Shared > Thumb Controller").
		SetPath(links.THUMB).
		SetHTMLHandler(thumb.NewThumbController(app).Handler)

	return []rtr.RouteInterface{
		adsTxt,
		files,
		flash,
		media,
		resources,
		theme,
		thumb,
	}
}
