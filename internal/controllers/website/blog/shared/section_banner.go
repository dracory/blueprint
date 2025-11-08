package shared

import (
	"project/internal/links"

	"github.com/dracory/hb"
)

func SectionBanner() *hb.Tag {
	style := hb.Style(`
.fill-success {
	fill: #0cbc87 !important;
}
.fill-orange {
	fill: #fd7e14 !important;
}
.fill-purple {
	fill: #6f42c1 !important;
}
	`)

	nav := hb.Nav().
		Class("breadcrumb mb-0").
		Attr("aria-label", "breadcrumb").
		Child(hb.OL().
			Class("breadcrumb mb-0").
			Child(hb.LI().
				Class("breadcrumb-item").
				Child(hb.A().
					Href(links.Website().Home()).
					HTML("Home"),
				)).
			Child(hb.LI().
				Class("breadcrumb-item active").
				Attr("aria-current", "page").
				Child(hb.A().
					Href(links.Website().Blog()).
					HTML("Blog"),
				),
			))

	section := hb.Section().
		Style("background:#1C1626;").
		Style("padding: 30px 0px;").
		Child(hb.Div().
			Class("container").
			Child(hb.Div().
				Class("row").
				Child(hb.Div().
					//HTML(decorationCross).
					Class("col-lg-10 mx-auto text-center").
					Style(`position: relative;`).
					Child(hb.I().Class("bi bi-crosshair").
						Style("color: magenta;").
						Style(`position: absolute; top: 0px; left: 0px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.I().Class("bi bi-asterisk").
						Style("color: magenta;").
						Style(`position: absolute; bottom: -10px; left: 100px;`).
						Style(`font-size: 30px; margin-right: 10px; transform: rotate(180deg);`)).
					Child(hb.I().Class("bi bi-star").
						Style("color: magenta;").
						Style(`position: absolute; top: 10px; right: 35px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.I().Class("bi bi-star").
						Style("color: magenta;").
						Style(`position: absolute; top: 15px; right: 0px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.I().Class("bi bi-star").
						Style("color: magenta;").
						Style(`position: absolute; top: 40px; right: 17px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.H1().Style("color:white;").HTML("Blog")).
					Child(hb.Div().
						Class("d-flex justify-content-center position-relative").
						Child(nav),
					),
				),
			),
		)

	return hb.Wrap().
		Child(style).
		Child(section)
}
