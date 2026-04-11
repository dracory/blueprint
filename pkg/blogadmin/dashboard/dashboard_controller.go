package dashboard

import (
	"context"
	"net/http"
	"strconv"

	"project/pkg/blogadmin/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/blogstore"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/uid"
)

// == CONTROLLER ==============================================================

type dashboardController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewDashboardController(registry registry.RegistryInterface) *dashboardController {
	return &dashboardController{registry: registry}
}

func (controller *dashboardController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, errorMessage, links.Admin().Blog(), 10)
	}

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Blog | Dashboard",
		Content: controller.page(data),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *dashboardController) page(data dashboardControllerData) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Blog",
			URL:  links.Admin().Blog(),
		},
		{
			Name: "Dashboard",
			URL:  "",
		},
	})

	// Add navigation tabs like CMS
	navTabs := controller.navTabs(data)

	title := hb.Heading1().
		HTML("Blog. Dashboard")

	content := hb.Wrap().
		Child(navTabs).
		Child(hb.BR()).
		Child(controller.dashboardCards(data))

	return layouts.AdminPage(
		title,
		breadcrumbs,
		content,
	)
}

func (controller *dashboardController) navTabs(data dashboardControllerData) hb.TagInterface {
	children := []hb.TagInterface{
		hb.Hyperlink().
			Class("text-decoration-none").
			HTML("Dashboard").
			Href(shared.NewLinks("/admin/blog").Dashboard()),
		hb.Hyperlink().
			Class("text-decoration-none").
			Child(
				hb.Wrap().
					Text("Posts ").
					Child(
						hb.Span().
							Class("badge bg-secondary ms-1").
							Text(strconv.FormatInt(data.postCount, 10)),
					),
			).
			Href(shared.NewLinks("/admin/blog").PostManager()),
	}

	// Only show Categories/Tags tabs if taxonomy is enabled
	if data.taxonomyEnabled {
		children = append(children,
			hb.Hyperlink().
				Class("text-decoration-none").
				Child(
					hb.Wrap().
						Text("Categories ").
						Child(
							hb.Span().
								Class("badge bg-secondary ms-1").
								Text(strconv.FormatInt(data.categoryCount, 10)),
						),
				).
				Href(shared.NewLinks("/admin/blog").CategoryManager()),
			hb.Hyperlink().
				Class("text-decoration-none").
				Child(
					hb.Wrap().
						Text("Tags ").
						Child(
							hb.Span().
								Class("badge bg-secondary ms-1").
								Text(strconv.FormatInt(data.tagCount, 10)),
						),
				).
				Href(shared.NewLinks("/admin/blog").TagManager()),
		)
	}

	return hb.Div().
		Class("card mb-4").
		Child(
			hb.Div().
				Class("card-body d-flex justify-content-center gap-4").
				Children(children),
		)
}

func (controller *dashboardController) dashboardCards(data dashboardControllerData) hb.TagInterface {
	cards := []hb.TagInterface{
		hb.Div().Class("col-md-4").Child(controller.cardPosts(data)),
	}

	// Only show category/tag cards if taxonomy is enabled
	if data.taxonomyEnabled {
		cards = append(cards,
			hb.Div().Class("col-md-4").Child(controller.cardCategories(data)),
			hb.Div().Class("col-md-4").Child(controller.cardTags(data)),
		)
	}

	return hb.Div().
		Class("row").
		Children(cards)
}

func (controller *dashboardController) cardPosts(data dashboardControllerData) hb.TagInterface {
	return hb.Hyperlink().
		Class("text-decoration-none").
		Href(shared.NewLinks("/admin/blog").PostManager()).
		Child(
			hb.Div().
				Class("card mb-4").
				Style("background-color: #a8d5ba; border: none;").
				Child(
					hb.Div().
						Class("card-body d-flex justify-content-between align-items-center").
						Children([]hb.TagInterface{
							hb.Div().
								Children([]hb.TagInterface{
									hb.Heading3().Class("mb-0").Text(strconv.FormatInt(data.postCount, 10)).Style("color: #2c5f2d;"),
									hb.Paragraph().Class("mb-0").Text("Total Posts").Style("color: #2c5f2d;"),
								}),
							hb.I().Class("bi bi-file-text fs-1").Style("color: rgba(44, 95, 45, 0.3);"),
						}),
				).
				Child(
					hb.Div().
						Class("card-footer bg-transparent border-0 text-center pb-3").
						Child(
							hb.Span().
								Class("small fw-medium").
								Text("More info ").
								Style("color: #2c5f2d;").
								Child(hb.I().Class("bi bi-arrow-right-circle")),
						),
				),
		)
}

func (controller *dashboardController) cardCategories(data dashboardControllerData) hb.TagInterface {
	return hb.Hyperlink().
		Class("text-decoration-none").
		Href(shared.NewLinks("/admin/blog").CategoryManager()).
		Child(
			hb.Div().
				Class("card mb-4").
				Style("background-color: #c5d5f5; border: none;").
				Child(
					hb.Div().
						Class("card-body d-flex justify-content-between align-items-center").
						Children([]hb.TagInterface{
							hb.Div().
								Children([]hb.TagInterface{
									hb.Heading3().Class("mb-0").Text(strconv.FormatInt(data.categoryCount, 10)).Style("color: #1a3a6e;"),
									hb.Paragraph().Class("mb-0").Text("Categories").Style("color: #1a3a6e;"),
								}),
							hb.I().Class("bi bi-folder fs-1").Style("color: rgba(26, 58, 110, 0.3);"),
						}),
				).
				Child(
					hb.Div().
						Class("card-footer bg-transparent border-0 text-center pb-3").
						Child(
							hb.Span().
								Class("small fw-medium").
								Text("More info ").
								Style("color: #1a3a6e;").
								Child(hb.I().Class("bi bi-arrow-right-circle")),
						),
				),
		)
}

func (controller *dashboardController) cardTags(data dashboardControllerData) hb.TagInterface {
	return hb.Hyperlink().
		Class("text-decoration-none").
		Href(shared.NewLinks("/admin/blog").TagManager()).
		Child(
			hb.Div().
				Class("card mb-4").
				Style("background-color: #f5e6c8; border: none;").
				Child(
					hb.Div().
						Class("card-body d-flex justify-content-between align-items-center").
						Children([]hb.TagInterface{
							hb.Div().
								Children([]hb.TagInterface{
									hb.Heading3().Class("mb-0").Text(strconv.FormatInt(data.tagCount, 10)).Style("color: #8b6914;"),
									hb.Paragraph().Class("mb-0").Text("Tags").Style("color: #8b6914;"),
								}),
							hb.I().Class("bi bi-tags fs-1").Style("color: rgba(139, 105, 20, 0.3);"),
						}),
				).
				Child(
					hb.Div().
						Class("card-footer bg-transparent border-0 text-center pb-3").
						Child(
							hb.Span().
								Class("small fw-medium").
								Text("More info ").
								Style("color: #8b6914;").
								Child(hb.I().Class("bi bi-arrow-right-circle")),
						),
				),
		)
}

func (controller *dashboardController) prepareData(r *http.Request) (data dashboardControllerData, errorMessage string) {
	ctx := r.Context()

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return data, "You are not logged in. Please login to continue."
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return data, "Blog store not available"
	}

	// Check if taxonomy is enabled by attempting to find a taxonomy
	_, err := blogStore.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_CATEGORY)
	if err != nil && err.Error() == "taxonomy is not enabled" {
		data.taxonomyEnabled = false
		data.taxonomyErrorMsg = "Categories and Tags are not available. To enable them, set TaxonomyEnabled to true in the blog store configuration."
	} else {
		data.taxonomyEnabled = true
		// Ensure taxonomies exist
		if err := controller.ensureTaxonomies(ctx, blogStore); err != nil {
			return data, "Failed to ensure taxonomies: " + err.Error()
		}
	}

	// Count posts
	postCount, err := blogStore.PostCount(ctx, blogstore.PostQueryOptions{})
	if err != nil {
		controller.registry.GetLogger().Error("blog dashboard: error counting posts", "error", err)
	}
	data.postCount = postCount

	// Get category taxonomy count only if taxonomy is enabled
	if data.taxonomyEnabled {
		categoryTaxonomy, err := blogStore.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_CATEGORY)
		if err == nil && categoryTaxonomy != nil {
			categoryCount, err := blogStore.TermCount(ctx, blogstore.TermQueryOptions{
				TaxonomyID: categoryTaxonomy.GetID(),
			})
			if err != nil {
				controller.registry.GetLogger().Error("blog dashboard: error counting categories", "error", err)
			}
			data.categoryCount = categoryCount
		}

		// Get tag taxonomy count
		tagTaxonomy, err := blogStore.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_TAG)
		if err == nil && tagTaxonomy != nil {
			tagCount, err := blogStore.TermCount(ctx, blogstore.TermQueryOptions{
				TaxonomyID: tagTaxonomy.GetID(),
			})
			if err != nil {
				controller.registry.GetLogger().Error("blog dashboard: error counting tags", "error", err)
			}
			data.tagCount = tagCount
		}
	}

	return data, ""
}

func (controller *dashboardController) ensureTaxonomies(ctx context.Context, store blogstore.StoreInterface) error {
	// Check if category taxonomy exists
	_, err := store.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_CATEGORY)
	if err != nil {
		controller.registry.GetLogger().Info("Creating category taxonomy")
		taxonomy := blogstore.NewTaxonomy()
		taxonomy.SetID(uid.HumanUid()[:8])
		taxonomy.SetName("Category")
		taxonomy.SetSlug(blogstore.TAXONOMY_CATEGORY)
		taxonomy.SetDescription("Blog post categories")
		if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
			return err
		}
	}

	// Check if tag taxonomy exists
	_, err = store.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_TAG)
	if err != nil {
		controller.registry.GetLogger().Info("Creating tag taxonomy")
		taxonomy := blogstore.NewTaxonomy()
		taxonomy.SetID(uid.HumanUid()[:8])
		taxonomy.SetName("Tag")
		taxonomy.SetSlug(blogstore.TAXONOMY_TAG)
		taxonomy.SetDescription("Blog post tags")
		if err := store.TaxonomyCreate(ctx, taxonomy); err != nil {
			return err
		}
	}

	return nil
}

type dashboardControllerData struct {
	postCount        int64
	categoryCount    int64
	tagCount         int64
	taxonomyEnabled  bool
	taxonomyErrorMsg string
}

