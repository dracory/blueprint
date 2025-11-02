package layouts

import "github.com/dracory/hb"

// AdminPage wraps admin page content with consistent spacing
func AdminPage(elements ...hb.TagInterface) *hb.Tag {
	wrapper := hb.Div().
		Class("container").
		Class("py-4")

	for _, el := range elements {
		if el == nil {
			continue
		}
		wrapper.Child(el)
	}

	return wrapper
}
