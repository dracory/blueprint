package layouts

import (
	"net/http"
	"project/app/widgets"
	"project/config"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/frontend"
	"github.com/gouniverse/hb"
)

func NewCmsLayout(options Options) *cmsLayout {
	layout := &cmsLayout{}
	layout.request = options.Request
	layout.title = options.Title // + " | " + config.AppName
	layout.content = options.Content
	layout.scriptURLs = options.ScriptURLs
	layout.scripts = options.Scripts
	layout.styleURLs = options.StyleURLs
	layout.styles = options.Styles
	return layout
}

type cmsLayout struct {
	request *http.Request
	// websiteSection string // i.e. Blog, Website
	title      string
	content    hb.TagInterface
	scriptURLs []string
	scripts    []string
	styleURLs  []string
	styles     []string
}

func (layout *cmsLayout) ToHTML() string {
	list := widgets.WidgetRegistry()

	shortcodes := []cmsstore.ShortcodeInterface{}
	for _, widget := range list {
		shortcodes = append(shortcodes, widget)
	}

	fe := frontend.New(frontend.Config{
		Store:  config.CmsStore,
		Logger: &config.Logger,
	})

	html, err := fe.TemplateRenderHtmlByID(layout.request, config.CmsUserTemplateID, struct {
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
		config.Logger.Error("At WebsiteLayout", "error", err.Error())
		return "Template error. Please try again later"
	}

	return html
}
