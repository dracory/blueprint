package layouts

import (
	"net/http"
	"project/internal/registry"

	"github.com/dracory/cmsstore/frontend"
	"github.com/dracory/hb"
)

func NewCmsLayout(
	registry registry.RegistryInterface,
	r *http.Request,
	options Options,
) *cmsLayout {
	layout := &cmsLayout{}
	layout.registry = registry
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
	registry   registry.RegistryInterface
	request    *http.Request
	title      string
	content    hb.TagInterface
	scriptURLs []string
	scripts    []string
	styleURLs  []string
	styles     []string
}

func (layout *cmsLayout) ToHTML() string {
	if layout.registry == nil {
		return "App is not initialized"
	}
	if layout.registry.GetConfig() == nil {
		return "Config is not initialized"
	}
	if !layout.registry.GetConfig().GetCmsStoreUsed() {
		return "Cms store is not used"
	}
	// if layout.registry.GetConfig().GetCmsStoreTemplateID() == "" {
	// 	return "Cms store template is not set"
	// }
	if layout.registry.GetCmsStore() == nil {
		return "Cms store is not initialized"
	}

	// list := widgets.WidgetRegistry(layout.registry)

	// shortcodes := []cmsstore.ShortcodeInterface{}
	// for _, widget := range list {
	// 	shortcodes = append(shortcodes, widget)
	// }

	fe := frontend.New(frontend.Config{
		Store:  layout.registry.GetCmsStore(),
		Logger: layout.registry.GetLogger(),
	})

	pageContent := ""

	for _, styleURL := range layout.styleURLs {
		pageContent += "<link rel='stylesheet' href='" + styleURL + "'>"
	}

	for _, style := range layout.styles {
		pageContent += "<style>" + style + "</style>"
	}

	pageContent += layout.content.ToHTML()

	for _, script := range layout.scripts {
		pageContent += "<script>" + script + "</script>"
	}

	for _, scriptURL := range layout.scriptURLs {
		pageContent += "<script src='" + scriptURL + "'></script>"
	}

	html, err := fe.TemplateRenderHtmlByID(
		layout.request,
		layout.registry.GetConfig().GetCmsStoreTemplateID(),
		struct {
			PageContent         string
			PageCanonicalURL    string
			PageMetaDescription string
			PageMetaKeywords    string
			PageMetaRobots      string
			PageTitle           string
			Language            string
		}{
			PageContent:         pageContent,
			PageCanonicalURL:    "",
			PageMetaDescription: "",
			PageMetaKeywords:    "",
			PageMetaRobots:      "",
			PageTitle:           layout.title,
			Language:            "en",
		})

	if err != nil {
		layout.registry.GetLogger().Error("At CmsLayout", "error", err.Error())
		return "Template error. Please try again later"
	}

	return html
}
