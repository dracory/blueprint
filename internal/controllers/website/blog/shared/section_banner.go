package shared

import (
	"project/internal/links"

	"github.com/dracory/hb"
)

func SectionBanner() *hb.Tag {
	nav := hb.Nav().
		Class("breadcrumb mb-0").
		Attr("aria-label", "breadcrumb").
		Child(hb.OL().
			Class("breadcrumb mb-0").
			Child(hb.LI().
				Class("breadcrumb-item").
				Child(hb.A().
					Href(links.Website().Home()).
					Class("text-secondary fw-bold text-decoration-none").
					HTML("Home"),
				)).
			Child(hb.LI().
				Class("breadcrumb-item active").
				Attr("aria-current", "page").
				Child(hb.A().
					Href(links.Website().Blog()).
					Class("text-secondary fw-bold text-decoration-none").
					HTML("Blog"),
				),
			))

	section := hb.Section().
		Style("padding: 40px 0 30px;").
		Child(hb.Div().
			Class("container").
			Child(hb.Div().
				Class("card rounded-5 p-4 p-md-5 overflow-hidden text-center mx-auto").
				Child(hb.Div().
					Class("card-body py-4").
					Child(hb.Div().
						Class("d-flex align-items-center justify-content-center mx-auto mb-3 rounded-4").
						Style("width: 56px; height: 56px; background-color: var(--cf-cyan);").
						Child(hb.I().Class("bi bi-journal-text fs-4"))).
					Child(hb.H1().
						Class("fw-black text-uppercase mb-3").
						Style("font-size: clamp(2rem, 5vw, 3rem); letter-spacing: -1px;").
						HTML("Blog")).
					Child(hb.P().
						Class("text-secondary fw-bold mb-3").
						HTML("Insights on email-native learning, productivity, and skill development")).
					Child(hb.Div().
						Class("d-flex justify-content-center").
						Child(nav)),
				),
			),
		)

	return section
}
