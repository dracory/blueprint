package layouts

import (
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/vaultstore"
)

// Options defines the options for the layout
type Options struct {
	Request        *http.Request
	AppName        string
	WebsiteSection string
	Title          string
	Content        hb.TagInterface
	ScriptURLs     []string
	Scripts        []string
	StyleURLs      []string
	Styles         []string
	VaultStore     vaultstore.StoreInterface
	VaultKey       string
}
