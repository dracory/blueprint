package aipostcontentupdate

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
)

func (c *Controller) handleFetchData(r *http.Request) string {
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
			c.app.GetLogger().Error("AI content editor: failed to load post", "error", err.Error())
		}
		return api.Error("Post not found").ToString()
	}
	if post == nil {
		return api.Error("Post not found").ToString()
	}

	blocks := MarkdownToBlocks(post.GetContent())
	blockData := make([]map[string]string, 0, len(blocks))
	for _, b := range blocks {
		blockData = append(blockData, map[string]string{
			"id":   b.ID,
			"type": string(b.Type),
			"text": b.Text,
		})
	}

	return api.SuccessWithData("Post loaded", map[string]any{
		"title":   post.GetTitle(),
		"summary": post.GetSummary(),
		"blocks":  blockData,
	}).ToString()
}
