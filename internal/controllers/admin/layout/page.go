package layout

import "github.com/dracory/hb"

// Page wraps admin page content with consistent spacing.
// Nil elements are skipped.
func Page(elements ...hb.TagInterface) *hb.Tag {
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
