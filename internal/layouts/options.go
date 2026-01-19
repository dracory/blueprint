package layouts

import (
	"github.com/dracory/hb"
)

// Options defines the options for the layout
type Options struct {
	AppName        string
	WebsiteSection string
	Title          string
	Content        hb.TagInterface
	ScriptURLs     []string
	Scripts        []string
	StyleURLs      []string
	Styles         []string
	CanonicalURL   string
}
