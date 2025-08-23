package widgets

import (
	"net/http"
	"project/internal/types"
	"project/pkg/testimonials"

	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

var _ Widget = (*testimonialsWidget)(nil) // verify it extends the interface

// == CONSTUCTOR ==============================================================

// NewPrintWidget creates a new instance of the print struct.
//
// Parameters:
//   - None
//
// Returns:
//   - *print - A pointer to the print struct
func NewTestimonialsWidget(app types.AppInterface) *testimonialsWidget {
	return &testimonialsWidget{app: app}
}

// == WIDGET ================================================================

// print is the struct that will be used to render the print shortcode.
//
// This shortcode is used to evaluate the result of the provided content
// and return it.
//
// It uses Otto as the engine.
type testimonialsWidget struct {
	app types.AppInterface
}

// == PUBLIC METHODS =========================================================

// Alias the shortcode alias to be used in the template.
func (t *testimonialsWidget) Alias() string {
	return "x-testimonials"
}

// Description a user-friendly description of the shortcode.
func (t *testimonialsWidget) Description() string {
	return "Renders the testimonials"
}

// Render implements the shortcode interface.
func (t *testimonialsWidget) Render(r *http.Request, content string, params map[string]string) string {
	testimonialList, err := testimonials.TestimonialList(t.app.GetEntityStore())

	if err != nil {
		return "Error: " + err.Error()
	}

	if len(testimonialList) < 1 {
		return "No testimonials found"
	}

	row := hb.Div().
		Class("row").
		Children(lo.Map(testimonialList, func(testimonial testimonials.Testimonial, index int) hb.TagInterface {


			stars := hb.Wrap().
				Child(hb.I().Class("bi bi-star-fill")).
				Child(hb.I().Class("bi bi-star-fill")).
				Child(hb.I().Class("bi bi-star-fill")).
				Child(hb.I().Class("bi bi-star-fill")).
				Child(hb.I().Class("bi bi-star-fill"))

			name := testimonial.FirstName()
			if testimonial.LastName() != "" {
				name += " "
				name += testimonial.LastName()[0:1]
				name += "."
			}

			card := hb.Div().
				Class("card mb-3").
				Child(hb.Div().
					Class("card-body").
					Child(hb.H5().
						Class("card-title").
						Child(stars)).
					Child(hb.P().
						Class("card-text").
						Child(hb.Text(testimonial.Quote()))).
					Child(hb.P().
						Class("card-text").
						Text(name).
						Text(",")).
					Text(" ").
					Text(testimonial.JobTitle()))

			return hb.Div().
				Class("col-sm-6").
				Child(card)
		}))

	return row.ToHTML()
}
