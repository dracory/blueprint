package admin

import (
	"net/http"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"project/pkg/webtheme"

	adminCmsStore "github.com/gouniverse/cmsstore/admin"
	"github.com/gouniverse/hb"
)

type cmsNewController struct {
	app types.AppInterface
}

func NewCmsNewController(app types.AppInterface) *cmsNewController {
	return &cmsNewController{app: app}
}

func (controller *cmsNewController) Handler(w http.ResponseWriter, r *http.Request) {
	admin, err := adminCmsStore.New(adminCmsStore.AdminOptions{
		Store:                  controller.app.GetCmsStore(),
		Logger:                 controller.app.GetLogger(),
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
			return layouts.NewAdminLayout(controller.app, r, layouts.Options{
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
		if logger := controller.app.GetLogger(); logger != nil {
			logger.Error("At admin > cmsNewController > Handler", "error", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	admin.Handle(w, r)
}
