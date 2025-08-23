package admin

import (
	"net/http"
	"project/internal/types"
)

type cmsController struct {
	app types.AppInterface
}

func NewCmsOldController(app types.AppInterface) *cmsController {
	return &cmsController{app: app}
}

func (controller cmsController) Handler(w http.ResponseWriter, r *http.Request) string {
	// config.Cms.SetFuncLayout(func(content string) string {
	// 	return layouts.NewAdminLayout(r, layouts.Options{
	// 		Title:   "CMS",
	// 		Content: hb.Raw(content),
	// 	}).ToHTML()
	// })

	// config.Cms.Router(w, r)

	return ""
}
