package aipostcontentupdate

import (
	"net/http"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/cdn"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

type Controller struct {
	registry registry.RegistryInterface
}

func NewController(registry registry.RegistryInterface) *Controller {
	return &Controller{registry: registry}
}

func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	if strings.TrimSpace(postID) == "" {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Post ID is required", shared.NewLinks().PostManager(), 10)
	}

	component := NewFormAiPostContentUpdate(c.registry)
	if component == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Failed to initialize AI content editor", shared.NewLinks().PostManager(), 10)
	}

	rendered := liveflux.SSR(component, map[string]string{
		"post_id": postID,
	})
	if rendered == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Error rendering AI content editor", shared.NewLinks().PostManager(), 10)
	}

	return layouts.NewAdminLayout(c.registry, r, layouts.Options{
		Title:   "Edit Post Content",
		Content: rendered,
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
			"https://cdn.jsdelivr.net/npm/sortablejs@1.15.0/Sortable.min.js",
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
	}).ToHTML()
}