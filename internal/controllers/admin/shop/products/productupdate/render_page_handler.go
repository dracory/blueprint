package productupdate

import (
	"net/http"

	"project/internal/controllers/admin/shop/products/productupdate/detailscomponent"
	"project/internal/controllers/admin/shop/products/productupdate/mediacomponent"
	metadatacomponent "project/internal/controllers/admin/shop/products/productupdate/metadatacomponent"
	"project/internal/controllers/admin/shop/products/productupdate/tagscomponent"
	"project/internal/controllers/admin/shop/shared"
	"project/internal/layouts"
	"project/internal/links"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

// handleRenderPage handles page rendering for GET requests
func (controller *productUpdateController) handleRenderPage(r *http.Request, product shopstore.ProductInterface, view string, productID string) string {
	pageContent := controller.renderPageContent(r, product, view, productID)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Edit Product | Shop",
		Content: hb.Wrap().HTML(pageContent),
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

// renderPageContent renders the product update page content with tabs
func (controller *productUpdateController) renderPageContent(r *http.Request, product shopstore.ProductInterface, view string, productID string) string {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Shop",
			URL:  shared.NewLinks().Home(map[string]string{}),
		},
		{
			Name: "Product Manager",
			URL:  shared.NewLinks().Products(map[string]string{}),
		},
		{
			Name: "Edit Product",
			URL:  shared.NewLinks().ProductUpdate(map[string]string{"product_id": productID}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.NewLinks().Products(map[string]string{}))

	heading := hb.Heading1().
		HTML("Shop. Product. Edit Product").
		Child(buttonCancel)

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "details", "active").
				Href(shared.NewLinks().ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "details",
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "media", "active").
				Href(shared.NewLinks().ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "media",
				})).
				HTML("Media"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "tags", "active").
				Href(shared.NewLinks().ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "tags",
				})).
				HTML("Tags"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "metadata", "active").
				Href(shared.NewLinks().ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "metadata",
				})).
				HTML("Metadata")))

	productTitle := hb.Heading2().
		Class("mb-3").
		Text("Product: ").
		Text(product.GetTitle())

	var body hb.TagInterface

	switch view {
	case "details":
		body = detailscomponent.Render(controller.registry, productID)
	case "media":
		body = mediacomponent.Render(controller.registry, productID)
	case "tags":
		body = tagscomponent.Render(controller.registry, productID)
	case "metadata":
		body = metadatacomponent.Render(controller.registry, productID)
	default:
		body = detailscomponent.Render(controller.registry, productID)
	}

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header d-flex justify-content-between align-items-center").
				Child(hb.Heading4().
					HTMLIf(view == "details", "Product Details").
					HTMLIf(view == "media", "Product Media").
					HTMLIf(view == "tags", "Product Tags").
					HTMLIf(view == "metadata", "Product Metadata").
					Style("margin-bottom:0;display:inline-block;")).
				ChildIf(view == "details",
					hb.Button().
						Type("button").
						Class("btn btn-primary btn-sm").
						ID("details-save-btn-top").
						Child(hb.I().Class("bi bi-save me-1")).
						Text("Save Details")).
				ChildIf(view == "metadata",
					hb.Button().
						Type("button").
						Class("btn btn-primary btn-sm").
						ID("metadata-save-btn-top").
						Child(hb.I().Class("bi bi-save me-1")).
						Text("Save Metadata")).
				ChildIf(view == "media",
					hb.Button().
						Type("button").
						Class("btn btn-primary btn-sm").
						ID("media-save-btn-top").
						Child(hb.I().Class("bi bi-save me-1")).
						Text("Save Media")).
				ChildIf(view == "tags",
					hb.Button().
						Type("button").
						Class("btn btn-primary btn-sm").
						ID("tags-save-btn-top").
						Child(hb.I().Class("bi bi-save me-1")).
						Text("Save Tags")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(body),
		).
		ChildIf(view == "details",
			hb.Div().
				Class("card-footer d-flex justify-content-end").
				Child(hb.Button().
					Type("button").
					Class("btn btn-primary").
					ID("details-save-btn-bottom").
					Child(hb.I().Class("bi bi-save me-2")).
					Text("Save Details"))).
		ChildIf(view == "metadata",
			hb.Div().
				Class("card-footer d-flex justify-content-end").
				Child(hb.Button().
					Type("button").
					Class("btn btn-primary").
					ID("metadata-save-btn-bottom").
					Child(hb.I().Class("bi bi-save me-2")).
					Text("Save Metadata"))).
		ChildIf(view == "media",
			hb.Div().
				Class("card-footer d-flex justify-content-end").
				Child(hb.Button().
					Type("button").
					Class("btn btn-primary").
					ID("media-save-btn-bottom").
					Child(hb.I().Class("bi bi-save me-2")).
					Text("Save Media"))).
		ChildIf(view == "tags",
			hb.Div().
				Class("card-footer d-flex justify-content-end").
				Child(hb.Button().
					Type("button").
					Class("btn btn-primary").
					ID("tags-save-btn-bottom").
					Child(hb.I().Class("bi bi-save me-2")).
					Text("Save Tags")))

	pageContent := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(productTitle).
		Child(tabs).
		Child(card)

	return pageContent.ToHTML()
}
