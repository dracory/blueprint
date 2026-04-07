package categories

import (
	"context"
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/shop/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
)

type categoryCreateController struct {
	registry registry.RegistryInterface
}

func NewCategoryCreateController(registry registry.RegistryInterface) *categoryCreateController {
	return &categoryCreateController{registry: registry}
}

func (controller *categoryCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
	if r.Method == http.MethodPost {
		return controller.handlePost(w, r)
	}

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Create Category | Shop",
		Content: hb.Wrap().HTML(controller.renderPage(r)),
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *categoryCreateController) renderPage(r *http.Request) string {
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
			Name: "Category Manager",
			URL:  shared.NewLinks().Categories(map[string]string{}),
		},
		{
			Name: "Create Category",
			URL:  shared.NewLinks().CategoryCreate(map[string]string{}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.NewLinks().Categories(map[string]string{}))

	heading := hb.Heading1().
		HTML("Shop. Category. Create Category").
		Child(buttonCancel)

	// Fetch categories for parent selection
	var categories []shopstore.CategoryInterface
	if controller.registry.GetShopStore() != nil {
		query := shopstore.NewCategoryQuery()
		categories, _ = controller.registry.GetShopStore().CategoryList(r.Context(), query)
	}

	// Build parent category options
	parentOptions := []hb.TagInterface{
		hb.Option().Value("").HTML("No parent (root category)"),
	}
	for _, cat := range categories {
		parentOptions = append(parentOptions, hb.Option().Value(cat.GetID()).HTML(cat.GetTitle()))
	}

	form := hb.Form().
		ID("FormCategoryCreate").
		Method(http.MethodPost).
		Action(shared.NewLinks().CategoryCreate(map[string]string{})).
		Children([]hb.TagInterface{
			hb.Div().Class("mb-3").Children([]hb.TagInterface{
				hb.Label().
					For("title").
					Class("form-label").
					HTML("Title"),
				hb.Input().
					Type("text").
					Class("form-control").
					ID("title").
					Name("title").
					Placeholder("Enter category title").
					Required(true),
			}),
			hb.Div().Class("mb-3").Children([]hb.TagInterface{
				hb.Label().
					For("description").
					Class("form-label").
					HTML("Description"),
				hb.Textarea().
					Class("form-control").
					ID("description").
					Name("description").
					Attr("rows", "3").
					Placeholder("Enter category description"),
			}),
			hb.Div().Class("mb-3").Children([]hb.TagInterface{
				hb.Label().
					For("status").
					Class("form-label").
					HTML("Status"),
				hb.Select().
					Class("form-select").
					ID("status").
					Name("status").
					Children([]hb.TagInterface{
						hb.Option().Value("draft").HTML("Draft"),
						hb.Option().Value("active").HTML("Active"),
						hb.Option().Value("inactive").HTML("Inactive"),
					}),
			}),
			hb.Div().Class("mb-3").Children([]hb.TagInterface{
				hb.Label().
					For("parent_id").
					Class("form-label").
					HTML("Parent Category"),
				hb.Select().
					Class("form-select").
					ID("parent_id").
					Name("parent_id").
					Children(parentOptions),
			}),
		})

	card := bs.Card().
		Class("mt-3").
		Children([]hb.TagInterface{
			bs.CardBody().Child(form),
			bs.CardFooter().
				Class("d-flex justify-content-end").
				Children([]hb.TagInterface{
					hb.Button().
						Type("submit").
						Class("btn btn-primary").
						Attr("form", "FormCategoryCreate").
						Child(hb.I().Class("bi bi-save me-2")).
						HTML("Create Category"),
				}),
		})

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(card)

	return content.ToHTML()
}

func (controller *categoryCreateController) handlePost(w http.ResponseWriter, r *http.Request) string {
	title := r.FormValue("title")
	description := r.FormValue("description")
	status := r.FormValue("status")
	parentID := r.FormValue("parent_id")

	if title == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Title is required", shared.NewLinks().CategoryCreate(map[string]string{}), 10)
	}

	category := shopstore.NewCategory()
	category.SetTitle(title)
	category.SetDescription(description)
	category.SetStatus(status)
	category.SetParentID(parentID)

	if err := controller.registry.GetShopStore().CategoryCreate(context.Background(), category); err != nil {
		slog.Error("Failed to create category", slog.String("error", err.Error()))
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Failed to create category", shared.NewLinks().CategoryCreate(map[string]string{}), 10)
	}

	return helpers.ToFlashSuccess(controller.registry.GetCacheStore(), w, r, "Category created successfully", shared.NewLinks().Categories(map[string]string{}), 10)
}
