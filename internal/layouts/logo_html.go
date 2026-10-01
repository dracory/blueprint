package layouts

import (
	_ "embed"
	"encoding/base64"

	"github.com/dracory/hb"
)

//go:embed logo.svg
var logoSVG string

// LogoHTML generates the HTML for the logo. The SVG is embedded as a data URI
// so it stays CSP-safe (img-src allows data:) without hotlinking an external
// domain.
func LogoHTML() string {
	src := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(logoSVG))
	img := hb.Image(src).
		Attr("alt", "Logo").
		Attr("height", "32").
		ToHTML()
	return img
}
