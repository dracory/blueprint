package aipostcontentupdate

import (
	"encoding/json"
	"net/http"
	"strings"

	"project/pkg/blogadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/blogstore"
)

func (c *Controller) handleSave(w http.ResponseWriter, r *http.Request) string {
	var reqBody struct {
		Action  string            `json:"action"`
		Title   string            `json:"title"`
		Summary string            `json:"summary"`
		Blocks  []map[string]string `json:"blocks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	postID := strings.TrimSpace(r.URL.Query().Get("post_id"))
	if postID == "" {
		postID = r.PostFormValue("post_id")
	}
	if postID == "" {
		return api.Error("Post ID is required").ToString()
	}

	if c.app.GetBlogStore() == nil {
		return api.Error("Blog store is not configured").ToString()
	}

	post, err := c.app.GetBlogStore().PostFindByID(r.Context(), postID)
	if err != nil {
		if c.app.GetLogger() != nil {
			c.app.GetLogger().Error("AI content editor: failed to load post for save", "error", err.Error())
		}
		return api.Error("Failed to load post").ToString()
	}
	if post == nil {
		return api.Error("Post not found").ToString()
	}

	// Reconstruct blocks from request
	blocks := make([]Block, 0, len(reqBody.Blocks))
	for _, b := range reqBody.Blocks {
		blocks = append(blocks, Block{
			ID:   b["id"],
			Type: BlockType(b["type"]),
			Text: b["text"],
		})
	}

	title := strings.TrimSpace(reqBody.Title)
	summary := strings.TrimSpace(reqBody.Summary)
	if title != "" {
		post.SetTitle(title)
	}
	if summary != "" {
		post.SetSummary(summary)
	}

	markdown := BlocksToMarkdown(blocks)
	post.SetContent(markdown)
	post.SetEditor(blogstore.POST_EDITOR_MARKDOWN)

	if err := c.app.GetBlogStore().PostUpdate(r.Context(), post); err != nil {
		if c.app.GetLogger() != nil {
			c.app.GetLogger().Error("AI content editor: failed to save post", "error", err.Error())
		}
		return api.Error("Failed to save post. Please try again later.").ToString()
	}

	if reqBody.Action == "save_close" {
		return api.SuccessWithData("Post saved successfully", map[string]any{
			"redirect_url": shared.NewLinks("/admin/blog").PostManager(),
		}).ToString()
	}

	return api.Success("Changes applied successfully").ToString()
}
