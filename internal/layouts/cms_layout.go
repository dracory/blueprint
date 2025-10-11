package layouts

import (
	"net/http"
	"project/internal/types"
	"project/internal/widgets"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/frontend"
	"github.com/dracory/hb"
)

func NewCmsLayout(
	app types.AppInterface,
	r *http.Request,
	options Options,
) *cmsLayout {
	layout := &cmsLayout{}
	layout.app = app
	layout.request = r
	layout.title = options.Title
	layout.content = options.Content
	layout.scriptURLs = options.ScriptURLs
	layout.scripts = options.Scripts
	layout.styleURLs = options.StyleURLs
	layout.styles = options.Styles
	return layout
}

type cmsLayout struct {
	app        types.AppInterface
	request    *http.Request
	title      string
	content    hb.TagInterface
	scriptURLs []string
	scripts    []string
	styleURLs  []string
	styles     []string
}

func (layout *cmsLayout) ToHTML() string {
	if layout.app == nil {
		return "App is not initialized"
	}
	if layout.app.GetConfig() == nil {
		return "Config is not initialized"
	}
	// if !layout.app.GetConfig().GetCmsStoreUsed() {
	// 	return "Cms store is not used"
	// }
	// if layout.app.GetConfig().GetCmsStoreTemplateID() == "" {
	// 	return "Cms store template is not set"
	// }
	if layout.app.GetCmsStore() == nil {
		return "Cms store is not initialized"
	}

	list := widgets.WidgetRegistry(layout.app)

	shortcodes := []cmsstore.ShortcodeInterface{}
	for _, widget := range list {
		shortcodes = append(shortcodes, widget)
	}

	fe := frontend.New(frontend.Config{
		Store:  layout.app.GetCmsStore(),
		Logger: layout.app.GetLogger(),
	})

	html, err := fe.TemplateRenderHtmlByID(
		layout.request,
		layout.app.GetConfig().GetCmsStoreTemplateID(),
		struct {
			PageContent         string
			PageCanonicalURL    string
			PageMetaDescription string
			PageMetaKeywords    string
			PageMetaRobots      string
			PageTitle           string
			Language            string
		}{
			PageContent:         layout.content.ToHTML(),
			PageCanonicalURL:    "",
			PageMetaDescription: "",
			PageMetaKeywords:    "",
			PageMetaRobots:      "",
			PageTitle:           layout.title,
			Language:            "en",
		})

	if err != nil {
		layout.app.GetLogger().Error("At WebsiteLayout", "error", err.Error())
		return "Template error. Please try again later"
	}

	return html
}
