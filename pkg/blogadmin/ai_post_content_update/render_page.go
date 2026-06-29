package aipostcontentupdate

import (
	"log/slog"
	"net/http"
	"strings"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/pkg/blogadmin/shared"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

func (c *Controller) renderPage(w http.ResponseWriter, r *http.Request) string {
	postID := strings.TrimSpace(r.URL.Query().Get("post_id"))
	if postID == "" {
		postID = r.PostFormValue("post_id")
	}
	if postID == "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, "Post ID is required", shared.NewLinks("/admin/blog").PostManager(), 10)
	}

	htmlContent, err := editorFiles.ReadFile("editor.html")
	if err != nil {
		slog.Error("Failed to read editor HTML", "error", err)
		return hb.Div().HTML("Error loading editor").ToHTML()
	}

	jsContent, err := editorFiles.ReadFile("editor.js")
	if err != nil {
		slog.Error("Failed to read editor JS", "error", err)
		return hb.Div().HTML("Error loading editor").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		window.postEditorPostId = '` + postID + `';
		window.postEditorBackUrl = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": postID}) + `';
		const urlPostEditorFetchData = '` + shared.NewLinks("/admin/blog").AiPostContentUpdate(map[string]string{"post_id": postID, "action": actionFetchData}) + `';
		const urlPostEditorRegenerate = '` + shared.NewLinks("/admin/blog").AiPostContentUpdate(map[string]string{"post_id": postID, "action": actionRegenerateBlock}) + `';
		const urlPostEditorSave = '` + shared.NewLinks("/admin/blog").AiPostContentUpdate(map[string]string{"post_id": postID, "action": actionSave}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:   "Edit Post Content",
		Content: vueContainer,
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
			"https://cdn.jsdelivr.net/npm/sortablejs@1.15.0/Sortable.min.js",
		},
	}).ToHTML()
}
