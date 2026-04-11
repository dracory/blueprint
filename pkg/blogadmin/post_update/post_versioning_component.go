package post_update

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	livefluxctl "project/internal/controllers/liveflux"
	"project/internal/registry"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
)

type postVersioningComponent struct {
	liveflux.Base

	registry registry.RegistryInterface

	PostID   string
	Post     blogstore.PostInterface
	Versions []blogstore.VersioningInterface

	FormErrorMessage   string
	FormSuccessMessage string
}

func NewPostVersioningComponent(registry registry.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&postVersioningComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*postVersioningComponent); ok {
		c.registry = registry
	}

	return inst
}

func (c *postVersioningComponent) GetKind() string {
	return "admin_blog_post_versioning"
}

func (c *postVersioningComponent) Mount(ctx context.Context, params map[string]string) error {
	if c.registry == nil {
		if registry, ok := ctx.Value(livefluxctl.AppContextKey).(registry.RegistryInterface); ok {
			c.registry = registry
		}
	}

	c.PostID = strings.TrimSpace(params["post_id"])
	if c.PostID == "" {
		c.FormErrorMessage = "Post ID is required"
		return nil
	}

	if c.registry == nil || c.registry.GetBlogStore() == nil {
		c.FormErrorMessage = "Blog store not available"
		return nil
	}

	// Load the post
	post, err := c.registry.GetBlogStore().PostFindByID(context.Background(), c.PostID)
	if err != nil {
		c.registry.GetLogger().Error("Error loading post for versioning component", "error", err.Error())
		c.FormErrorMessage = "Error loading post"
		return nil
	}

	if post == nil {
		c.FormErrorMessage = "Post not found"
		return nil
	}

	c.Post = post

	// Load versions
	if err := c.loadVersions(); err != nil {
		c.registry.GetLogger().Error("Error loading versions", "error", err.Error())
		c.FormErrorMessage = "Error loading versions"
		return nil
	}

	return nil
}

func (c *postVersioningComponent) loadVersions() error {
	if !c.registry.GetBlogStore().VersioningEnabled() {
		return nil
	}

	versions, err := c.registry.GetBlogStore().VersioningList(context.Background(), blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(c.PostID).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC).
		SetLimit(50)) // Limit to last 50 versions

	if err != nil {
		return err
	}

	c.Versions = versions
	return nil
}

func (c *postVersioningComponent) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "restore":
		versionID := strings.TrimSpace(data.Get("version_id"))
		if versionID == "" {
			c.FormErrorMessage = "Version ID is required"
			return nil
		}

		if c.registry == nil || c.registry.GetBlogStore() == nil {
			c.FormErrorMessage = "Blog store not available"
			return nil
		}

		// Get the version to restore
		version, err := c.registry.GetBlogStore().VersioningFindByID(context.Background(), versionID)
		if err != nil || version == nil {
			c.FormErrorMessage = "Version not found"
			return nil
		}

		// Get the current post
		post, err := c.registry.GetBlogStore().PostFindByID(context.Background(), c.PostID)
		if err != nil || post == nil {
			c.FormErrorMessage = "Post not found"
			return nil
		}

		// Restore the post from version content
		if err := post.UnmarshalFromVersioning(version.Content()); err != nil {
			c.registry.GetLogger().Error("Error unmarshaling version content", "error", err.Error())
			c.FormErrorMessage = "Error restoring version"
			return nil
		}

		// Update the post
		if err := c.registry.GetBlogStore().PostUpdate(context.Background(), post); err != nil {
			c.registry.GetLogger().Error("Error updating post from version", "error", err.Error())
			c.FormErrorMessage = "Error restoring version"
			return nil
		}

		// Create a new version for the restoration
		if err := createPostVersioning(context.Background(), c.registry, post); err != nil {
			c.registry.GetLogger().Error("Error creating post versioning after restore", "error", err.Error())
		}

		// Reload post and versions
		c.Post = post
		if err := c.loadVersions(); err != nil {
			c.registry.GetLogger().Error("Error reloading versions after restore", "error", err.Error())
		}

		c.FormErrorMessage = ""
		c.FormSuccessMessage = "Version restored successfully"
	}

	return nil
}

func (c *postVersioningComponent) Render(ctx context.Context) hb.TagInterface {
	if !c.registry.GetBlogStore().VersioningEnabled() {
		return c.Root(hb.Div().Class("alert alert-info").Text("Versioning is not enabled"))
	}

	if len(c.Versions) == 0 {
		return c.Root(hb.Div().Class("alert alert-info").Text("No versions available for this post"))
	}

	// Create version list
	versionRows := []hb.TagInterface{}

	for i, version := range c.Versions {
		// Parse version content to get post data
		var postData map[string]interface{}
		if err := json.Unmarshal([]byte(version.Content()), &postData); err != nil {
			continue
		}

		title, _ := postData["title"].(string)
		summary, _ := postData["summary"].(string)

		if title == "" {
			title = "Untitled"
		}

		// Format the creation time
		createdAtStr := version.CreatedAt()
		timeStr := formatDateTime(createdAtStr)

		// Create restore button
		restoreBtn := hb.Button().
			Type("button").
			Class("btn btn-sm btn-outline-primary js-restore-version").
			Attr("data-version-id", version.ID()).
			Child(hb.I().Class("bi bi-arrow-clockwise me-1")).
			Child(hb.Span().Text("Restore"))

		// Create view details button
		viewBtn := hb.Button().
			Type("button").
			Class("btn btn-sm btn-outline-info me-2 js-view-version").
			Attr("data-version-id", version.ID()).
			Attr("data-bs-toggle", "modal").
			Attr("data-bs-target", "#versionModal").
			Child(hb.I().Class("bi bi-eye me-1")).
			Child(hb.Span().Text("View"))

		// Hidden action button for liveflux
		actionBtn := hb.Button().
			Type("button").
			Class("d-none js-restore-version-action").
			Attr(liveflux.DataFluxAction, "restore").
			Attr(liveflux.DataFluxTargetKind, c.GetKind()).
			Attr(liveflux.DataFluxTargetID, c.GetID()).
			Child(hb.Span().Text("Restore"))

		// Add version_id as hidden input for the restore action
		restoreForm := hb.Form().
			Class("d-inline").
			Attr("method", "post").
			Child(hb.Input().
				Type("hidden").
				Name("version_id").
				Value(version.ID())).
			Child(viewBtn).
			Child(restoreBtn).
			Child(actionBtn)

		// Version content preview
		contentPreview := ""
		if summary != "" {
			contentPreview = summary
			if len(contentPreview) > 100 {
				contentPreview = contentPreview[:100] + "..."
			}
		}

		// Get changes compared to previous version
		var changes []string
		if i < len(c.Versions)-1 {
			// Get previous version data
			var previousPostData map[string]interface{}
			if err := json.Unmarshal([]byte(c.Versions[i+1].Content()), &previousPostData); err == nil {
				changes = getVersionChanges(postData, previousPostData)
			}
		} else {
			changes = []string{"Initial version"}
		}

		// Create changes badges
		changesBadges := hb.Span()
		for _, change := range changes {
			changesBadges = changesBadges.Child(
				hb.Span().
					Class("badge bg-secondary me-1").
					Text(change),
			)
		}

		// Create version row
		row := hb.Tr().
			Child(hb.Td().Class("text-nowrap").Text(timeStr)).
			Child(hb.Td().
				Child(hb.Div().
					Class("fw-bold").
					Text(title)).
				Child(hb.Div().
					Class("text-muted small").
					Text(contentPreview)).
				Child(hb.Div().
					Class("mt-1").
					Child(changesBadges))).
			Child(hb.Td().Class("text-end").Child(restoreForm))

		// Highlight current version (first one)
		if i == 0 {
			row = row.Class("table-success")
		}

		versionRows = append(versionRows, row)
	}

	// Create version table
	table := hb.Table().
		Class("table table-striped table-hover").
		Child(hb.Thead().
			Child(hb.Tr().
				Child(hb.Th().Style("width: 180px;").Text("Date/Time")).
				Child(hb.Th().Text("Version")).
				Child(hb.Th().Style("width: 200px;").Class("text-end").Text("Actions"))))

	// Add tbody with rows
	tbody := hb.Tbody()
	for _, row := range versionRows {
		tbody = tbody.Child(row)
	}
	table = table.Child(tbody)

	// Create version details modal
	versionModal := hb.Div().
		Class("modal fade").
		ID("versionModal").
		Attr("tabindex", "-1").
		Child(hb.Div().
			Class("modal-dialog modal-lg").
			Child(hb.Div().
				Class("modal-content").
				Child(hb.Div().
					Class("modal-header").
					Child(hb.H5().
						Class("modal-title").
						Text("Version Details")).
					Child(hb.Button().
						Type("button").
						Class("btn-close").
						Attr("data-bs-dismiss", "modal")).
					Child(hb.Button().
						Type("button").
						Class("btn-close").
						Attr("data-bs-dismiss", "modal"))).
				Child(hb.Div().
					Class("modal-body").
					ID("versionModalBody").
					Text("Loading...")).
				Child(hb.Div().
					Class("modal-footer").
					Child(hb.Button().
						Type("button").
						Class("btn btn-secondary").
						Attr("data-bs-dismiss", "modal").
						Text("Close")))))

	// Create help text
	helpText := hb.Div().
		Class("alert alert-info mb-3").
		Child(hb.I().Class("bi bi-info-circle me-2")).
		Child(hb.Span().Text("Versions are automatically created when you save a post. The most recent version is highlighted in green. Click 'View' to see detailed changes or 'Restore' to revert to a previous version."))

	// Add SweetAlert confirmation script for restore and view functionality
	confirmScript := hb.Script(`
// Store version data globally
window.versionData = ` + fmt.Sprintf(`%s`, func() string {
		versionsData := map[string]map[string]interface{}{}
		for _, version := range c.Versions {
			var postData map[string]interface{}
			if err := json.Unmarshal([]byte(version.Content()), &postData); err == nil {
				versionsData[version.ID()] = postData
			}
		}
		data, _ := json.Marshal(versionsData)
		return string(data)
	}()) + `;

document.addEventListener('click', function(e) {
	// Handle restore button clicks
	const restoreBtn = e.target.closest('.js-restore-version');
	if (restoreBtn && typeof Swal !== 'undefined') {
		e.preventDefault();
		
		const versionId = restoreBtn.getAttribute('data-version-id');
		const actionBtn = document.querySelector('.js-restore-version-action');
		const hiddenInput = actionBtn.closest('form').querySelector('input[name="version_id"]');
		
		if (hiddenInput) {
			hiddenInput.value = versionId;
		}
		
		Swal.fire({
			title: 'Restore Version',
			text: 'Are you sure you want to restore this version? This will replace the current content.',
			icon: 'warning',
			showCancelButton: true,
			confirmButtonText: 'Yes, restore',
			cancelButtonText: 'Cancel'
		}).then((result) => {
			if (result.isConfirmed && actionBtn) {
				actionBtn.click();
			}
		});
	}
	
	// Handle view button clicks
	const viewBtn = e.target.closest('.js-view-version');
	if (viewBtn) {
		e.preventDefault();
		
		const versionId = viewBtn.getAttribute('data-version-id');
		const modalBody = document.getElementById('versionModalBody');
		
		if (modalBody && window.versionData[versionId]) {
			const version = window.versionData[versionId];
			
			let content = '<div class="row">';
			content += '<div class="col-md-6"><strong>Title:</strong></div>';
			content += '<div class="col-md-6">' + (version.title || 'Untitled') + '</div>';
			content += '</div>';
			
			content += '<div class="row mt-2">';
			content += '<div class="col-md-6"><strong>Status:</strong></div>';
			content += '<div class="col-md-6"><span class="badge bg-' + (version.status === 'published' ? 'success' : 'secondary') + '">' + (version.status || 'draft') + '</span></div>';
			content += '</div>';
			
			content += '<div class="row mt-2">';
			content += '<div class="col-md-6"><strong>Featured:</strong></div>';
			content += '<div class="col-md-6"><span class="badge bg-' + (version.featured === 'yes' ? 'warning' : 'secondary') + '">' + (version.featured || 'no') + '</span></div>';
			content += '</div>';
			
			if (version.summary) {
				content += '<div class="row mt-3">';
				content += '<div class="col-12"><strong>Summary:</strong></div>';
				content += '<div class="col-12"><p class="text-muted">' + version.summary + '</p></div>';
				content += '</div>';
			}
			
			if (version.content) {
				content += '<div class="row mt-3">';
				content += '<div class="col-12"><strong>Content Preview:</strong></div>';
				content += '<div class="col-12"><div class="border rounded p-2 bg-light" style="max-height: 300px; overflow-y: auto;"><pre class="mb-0 small">' + version.content.substring(0, 1000) + (version.content.length > 1000 ? '...' : '') + '</pre></div></div>';
				content += '</div>';
			}
			
			modalBody.innerHTML = content;
		}
	}
});`).ToHTML()

	content := hb.Div().
		Child(helpText).
		Child(table).
		Child(versionModal).
		Child(hb.Raw(confirmScript))

	// Add messages
	if c.FormErrorMessage != "" {
		content = content.Child(hb.Div().Class("alert alert-danger mt-3").Text(c.FormErrorMessage))
	}

	if c.FormSuccessMessage != "" {
		content = content.Child(hb.Div().Class("alert alert-success mt-3").Text(c.FormSuccessMessage))
	}

	return c.Root(content)
}

func init() {
	if err := liveflux.Register(&postVersioningComponent{}); err != nil {
		log.Printf("Failed to register postVersioningComponent: %v", err)
	}
}

// formatDateTime formats a datetime string to a user-friendly format
func formatDateTime(datetimeStr string) string {
	if datetimeStr == "" {
		return "Unknown"
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", datetimeStr)
	if err != nil {
		// Try other common formats
		parsedTime, err = time.Parse(time.RFC3339, datetimeStr)
		if err != nil {
			return "Unknown"
		}
	}

	now := time.Now()
	diff := now.Sub(parsedTime)
	days := int(diff.Hours() / 24)

	if days == 0 {
		return parsedTime.Format("Today, 3:04 PM")
	} else if days == 1 {
		return parsedTime.Format("Yesterday, 3:04 PM")
	} else if days < 7 {
		return parsedTime.Format("Monday, 3:04 PM")
	} else {
		return parsedTime.Format("Jan 2, 2006 3:04 PM")
	}
}

// getVersionChanges compares two versions and returns a summary of changes
func getVersionChanges(currentVersion, previousVersion map[string]interface{}) []string {
	changes := []string{}

	// Compare title
	if currentVersion["title"] != previousVersion["title"] {
		changes = append(changes, "Title changed")
	}

	// Compare summary
	if currentVersion["summary"] != previousVersion["summary"] {
		changes = append(changes, "Summary changed")
	}

	// Compare content
	currentContent, _ := currentVersion["content"].(string)
	previousContent, _ := previousVersion["content"].(string)

	if currentContent != previousContent {
		// Simple content length comparison for change detection
		currentLen := len(currentContent)
		previousLen := len(previousContent)

		if currentLen > previousLen {
			changes = append(changes, fmt.Sprintf("Content expanded (+%d chars)", currentLen-previousLen))
		} else if currentLen < previousLen {
			changes = append(changes, fmt.Sprintf("Content shortened (%d chars)", currentLen))
		} else {
			changes = append(changes, "Content modified")
		}
	}

	// Compare status
	if currentVersion["status"] != previousVersion["status"] {
		changes = append(changes, fmt.Sprintf("Status changed to %s", currentVersion["status"]))
	}

	// Compare featured
	if currentVersion["featured"] != previousVersion["featured"] {
		changes = append(changes, fmt.Sprintf("Featured: %s", currentVersion["featured"]))
	}

	if len(changes) == 0 {
		changes = append(changes, "Minor changes")
	}

	return changes
}

