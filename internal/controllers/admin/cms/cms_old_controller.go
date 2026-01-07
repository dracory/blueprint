package admin

import (
	"net/http"
	"project/internal/registry"
)

type cmsController struct {
	registry registry.RegistryInterface
}

func NewCmsOldController(registry registry.RegistryInterface) *cmsController {
	return &cmsController{registry: registry}
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
