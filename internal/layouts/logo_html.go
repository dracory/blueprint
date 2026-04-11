package layouts

import "github.com/dracory/hb"

// LogoHTML generates the HTML for the logo
func LogoHTML() string {
	img := hb.Image("https://dracory.com/assets/images/logo.png").ToHTML()
	return img
}
