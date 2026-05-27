// Package shopadmin provides a shop admin interface following the folder-per-controller pattern.
// Each controller is in its own subfolder and handles its own views and AJAX data.
// This structure allows for future migration to a standalone external package.
package shopadmin

import (
	"net/http"
	"strings"

	"project/internal/registry"
)

// AdminOptions contains all dependencies and configuration for the shop admin
type AdminOptions struct {
	// Registry provides access to all stores and services
	Registry registry.RegistryInterface

	// AdminHomeURL is the URL for the admin home page
	AdminHomeURL string

	// ShopAdminURL is the base URL for the shop admin (e.g., "/admin/shop")
	ShopAdminURL string

	// AuthUserID returns the authenticated user ID from the request
	AuthUserID func(r *http.Request) string

	// FileManagerURL is the URL for the file manager (e.g., "/admin/files")
	FileManagerURL string
}

// AdminInterface defines the interface for the shop admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	opts AdminOptions
}

// New creates a new shop admin instance
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Registry == nil {
		return nil, ErrRegistryRequired
	}

	// Set defaults
	if opts.ShopAdminURL == "" {
		opts.ShopAdminURL = "/admin/shop"
	}

	return &admin{opts: opts}, nil
}

// Handle processes all shop admin requests
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if a.opts.AuthUserID != nil && a.opts.AuthUserID(r) == "" {
		http.Redirect(w, r, a.opts.AdminHomeURL, http.StatusSeeOther)
		return
	}

	// Use Routes() with registry to handle the request
	if a.opts.Registry == nil {
		http.Error(w, "Registry not configured", http.StatusInternalServerError)
		return
	}

	// Get routes and find matching one
	routes, err := Routes(a.opts.Registry, a.opts)
	if err != nil {
		http.Error(w, "Failed to load routes", http.StatusInternalServerError)
		return
	}

	// Find matching route by path
	for _, route := range routes {
		if strings.HasPrefix(r.URL.Path, route.GetPath()) {
			// Execute the route's handler
			if handler := route.GetHandler(); handler != nil {
				handler(w, r)
				return
			}
			if htmlHandler := route.GetHTMLHandler(); htmlHandler != nil {
				htmlHandler(w, r)
				return
			}
		}
	}

	// No route matched - redirect to home
	http.Redirect(w, r, a.opts.ShopAdminURL+"?controller=home", http.StatusSeeOther)
}
