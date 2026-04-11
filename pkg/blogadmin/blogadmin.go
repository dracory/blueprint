// Package blogadmin provides a blog admin interface following the folder-per-controller pattern.
// Each controller is in its own subfolder and handles its own views and AJAX data.
// This structure allows for future migration to a standalone external package.
package blogadmin

import (
	"net/http"
	"strings"

	"github.com/dracory/blogstore"

	"project/internal/registry"
)

// LLMEngineInterface defines the interface for LLM operations
type LLMEngineInterface interface {
	GenerateText(systemPrompt string, userPrompt string) (string, error)
}

// AdminOptions contains all dependencies and configuration for the blog admin
type AdminOptions struct {
	// Store is the blog store instance
	Store blogstore.StoreInterface

	// AdminHomeURL is the URL for the admin home page
	AdminHomeURL string

	// BlogAdminURL is the base URL for the blog admin (e.g., "/admin/blog")
	BlogAdminURL string

	// AuthUserID returns the authenticated user ID from the request
	AuthUserID func(r *http.Request) string

	// LLMEngine is the AI/LLM engine for generating content (optional)
	LLMEngine LLMEngineInterface

	// BlogTopic is the topic for AI-generated content (optional)
	BlogTopic string

	// Registry provides access to all stores and services
	Registry registry.RegistryInterface
}

// AdminInterface defines the interface for the blog admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	opts AdminOptions
}

// New creates a new blog admin instance
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Store == nil {
		return nil, ErrStoreRequired
	}

	// Set defaults
	if opts.BlogAdminURL == "" {
		opts.BlogAdminURL = "/admin/blog"
	}

	return &admin{opts: opts}, nil
}

// Handle processes all blog admin requests
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
	routes := Routes(a.opts.Registry)

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

	// No route matched - redirect to dashboard
	http.Redirect(w, r, a.opts.BlogAdminURL+"?controller=dashboard", http.StatusSeeOther)
}

