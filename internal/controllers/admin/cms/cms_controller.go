package admin

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/webtheme"

	adminCmsStore "github.com/dracory/cmsstore/admin"
	"github.com/dracory/hb"
)

type cmsNewController struct {
	registry registry.RegistryInterface
}

func NewCmsNewController(registry registry.RegistryInterface) *cmsNewController {
	return &cmsNewController{registry: registry}
}

func (controller *cmsNewController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := adminCmsStore.New(adminCmsStore.AdminOptions{
		Store:                  controller.registry.GetCmsStore(),
		Logger:                 controller.registry.GetLogger(),
		BlockEditorDefinitions: webtheme.BlockEditorDefinitions(),
		// BlockEditorRenderer: func(blocks []ui.BlockInterface) string {
		// 	return webtheme.New(blocks).ToHtml()
		// },
		FuncLayout: func(pageTitle string, pageContent string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
				Title:      pageTitle + " | CMS (NEW)",
				Content:    hb.Raw(pageContent),
				ScriptURLs: options.ScriptURLs,
				StyleURLs:  options.StyleURLs,
				Scripts:    options.Scripts,
				Styles:     options.Styles,
			}).ToHTML()
		},
		AdminHomeURL: links.Admin().Home(),
	})

	if err != nil {
		if logger := controller.registry.GetLogger(); logger != nil {
			logger.Error("At admin > cmsNewController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			if logger := controller.registry.GetLogger(); logger != nil {
				logger.Error("At admin > cmsNewController > Handler", "write_error", writeErr.Error())
			}
		}

		return
	}

	admin.Handle(w, r)
}
