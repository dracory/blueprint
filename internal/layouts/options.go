package layouts

import (
	"github.com/dracory/hb"
)

// Options defines the options for the layout
type Options struct {
	// Request        *http.Request
	AppName        string
	WebsiteSection string
	Title          string
	Content        hb.TagInterface
	ScriptURLs     []string
	Scripts        []string
	StyleURLs      []string
	Styles         []string
}
