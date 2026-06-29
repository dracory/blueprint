package layouts

import (
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/links"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// userLayoutNavbar builds the neo-brutalist navbar for the user layout.
func userLayoutNavbar(app app.AppInterface, r *http.Request) hb.TagInterface {
	authUser := helpers.GetAuthUser(r)
	firstName := ""
	if authUser != nil {
		fn, _, _ := userDisplayNames(app, r, authUser, app.GetConfig().GetVaultStoreKey())
		firstName = fn
	}

	mainMenuItems := userLayoutMainMenuItems(authUser)
	userMenuItems := userLayoutUserMenuItems(authUser)

	// Brand link
	brandLink := hb.Hyperlink().
		Href(links.Website().Home()).
		Class("navbar-brand d-flex align-items-center gap-3").
		HTML(LogoHTML())

	// Navbar collapse toggle for mobile
	toggleButton := hb.Button().
		Type("button").
		Class("navbar-toggler").
		Attr("data-bs-toggle", "collapse").
		Attr("data-bs-target", "#userNavbarCollapse").
		Attr("aria-controls", "userNavbarCollapse").
		Attr("aria-expanded", "false").
		Attr("aria-label", "Toggle navigation").
		Child(hb.Span().Class("navbar-toggler-icon"))

	// Main nav links (visible on desktop)
	var navLinks []hb.TagInterface
	for _, item := range mainMenuItems {
		link := hb.Hyperlink().
			Href(item.URL).
			Class("nav-link fw-bold text-uppercase tracking-wider small")
		if item.Target != "" {
			link.Attr("target", item.Target)
		}
		link.HTML(lo.Ternary(item.Icon != "", item.Icon+" ", "") + item.Title)
		navLinks = append(navLinks, hb.LI().Class("nav-item").Child(link))
	}

	mainNav := hb.UL().
		Class("navbar-nav me-auto mb-2 mb-lg-0").
		Children(navLinks)

	// User dropdown
	var dropdownItems []hb.TagInterface
	for _, item := range userMenuItems {
		link := hb.Hyperlink().
			Href(item.URL).
			Class("dropdown-item fw-bold py-2")
		if item.Target != "" {
			link.Attr("target", item.Target)
		}
		if item.Title == "Logout" {
			link.Class("dropdown-item fw-bold py-2 text-danger")
		}
		link.HTML(lo.Ternary(item.Icon != "", item.Icon+" ", "") + item.Title)
		dropdownItems = append(dropdownItems, hb.LI().Child(link))
	}

	// Divider before logout (which is the last item)
	if len(dropdownItems) > 1 {
		lastIdx := len(dropdownItems) - 1
		dropdownItems = append(
			dropdownItems[:lastIdx],
			append([]hb.TagInterface{hb.LI().Child(hb.HR().Class("dropdown-divider mx-2"))}, dropdownItems[lastIdx:]...)...,
		)
	}

	dropdownMenu := hb.UL().
		Class("dropdown-menu dropdown-menu-end rounded-4 mt-2").
		Children(dropdownItems)

	dropdownToggle := hb.Button().
		Type("button").
		Class("btn btn-dark rounded-4 px-3 py-2 fw-black text-uppercase tracking-wider small dropdown-toggle").
		Attr("data-bs-toggle", "dropdown").
		Attr("aria-expanded", "false").
		HTML(`<i class="bi bi-person-circle me-2"></i>` + lo.Ternary(firstName != "", firstName, "Account"))

	userDropdown := hb.Div().
		Class("dropdown").
		Child(dropdownToggle).
		Child(dropdownMenu)

	// Theme toggle button
	themeToggle := hb.Button().
		Type("button").
		ID("themeToggle").
		Class("btn btn-outline-dark rounded-4 px-3 py-2 fw-black").
		Attr("title", "Toggle theme").
		Child(hb.I().Class("bi bi-moon").ID("themeIcon"))

	// Right side controls
	var rightSideChildren []hb.TagInterface
	rightSideChildren = append(rightSideChildren, mainNav)
	if authUser != nil {
		rightSideChildren = append(rightSideChildren, userDropdown)
	}
	rightSideChildren = append(rightSideChildren, themeToggle)

	rightSide := hb.Div().
		Class("d-flex align-items-center gap-3 ms-auto").
		Children(rightSideChildren)

	// Collapsible content
	collapse := hb.Div().
		ID("userNavbarCollapse").
		Class("collapse navbar-collapse").
		Child(rightSide)

	// Full navbar
	navbar := hb.Nav().
		Class("navbar navbar-expand-lg py-4").
		Style("max-width: 1200px; margin: 0 auto; width: 100%;").
		Child(hb.Div().Class("container").Child(brandLink).Child(toggleButton).Child(collapse))

	return hb.Section().ID("SectionNavbar").Child(navbar)
}
