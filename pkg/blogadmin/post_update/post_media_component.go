package post_update

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"project/internal/app"
	"project/pkg/blogadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

type postMediaComponent struct {
	app    app.AppInterface
	postID string
}

func newPostMediaComponent(app app.AppInterface) *postMediaComponent {
	return &postMediaComponent{app: app}
}

func (c *postMediaComponent) Render(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_media.html")
	if err != nil {
		slog.Error("Failed to read post media HTML template", "error", err)
		return hb.Div().HTML("Error loading media component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_media.js")
	if err != nil {
		slog.Error("Failed to read post media JavaScript file", "error", err)
		return hb.Div().HTML("Error loading media component")
	}

	initScript := `
		const postId = '` + post.GetID() + `';
		const urlMediaLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-media"}) + `';
		const urlMediaUpload = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "upload-media"}) + `';
		const urlMediaSave = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "save-media"}) + `';
		const urlMediaDelete = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"action": "delete-media"}) + `';
		const urlMediaAdd = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "add-media"}) + `';
	`

	return hb.Div().
		Child(hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")).
		Child(hb.Wrap().HTML(string(htmlContent))).
		Child(hb.Script(initScript)).
		Child(hb.Script(string(jsContent)))
}

func (c *postMediaComponent) HandleLoad(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("post_id is required").ToString()
	}

	files, err := c.app.GetBlogStore().MediaListByEntityID(r.Context(), postID)
	if err != nil {
		c.app.GetLogger().Error("Failed to load post files", "error", err)
		return api.Error("Failed to load files").ToString()
	}

	fileList := []map[string]any{}
	for _, f := range files {
		fileList = append(fileList, map[string]any{
			"id":        f.GetID(),
			"name":      f.GetTitle(),
			"url":       f.GetURL(),
			"type":      f.GetType(),
			"size":      f.GetSize(),
			"extension": f.GetExtension(),
			"sequence":  f.GetSequence(),
		})
	}

	return api.SuccessWithData("Files loaded successfully", map[string]any{
		"files": fileList,
	}).ToString()
}

func (c *postMediaComponent) HandleUpload(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("post_id is required").ToString()
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		return api.Error("Failed to parse upload: " + err.Error()).ToString()
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		files = r.MultipartForm.File["upload_file"]
	}
	if len(files) == 0 {
		return api.Error("No files uploaded").ToString()
	}

	existingFiles, _ := c.app.GetBlogStore().MediaListByEntityID(r.Context(), postID)
	startSequence := len(existingFiles)

	uploaded := []map[string]any{}

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return api.Error("Failed to open file: " + err.Error()).ToString()
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return api.Error("Failed to read file: " + err.Error()).ToString()
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		dataURI := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data)

		postFile := blogstore.NewMedia().
			SetEntityID(postID).
			SetTitle(fileHeader.Filename).
			SetURL(dataURI).
			SetType(contentType).
			SetSize(strconv.FormatInt(fileHeader.Size, 10)).
			SetExtension(ext).
			SetSequence(startSequence + i)

		if err := c.app.GetBlogStore().MediaCreate(r.Context(), postFile); err != nil {
			return api.Error("Failed to save file record: " + err.Error()).ToString()
		}

		uploaded = append(uploaded, map[string]any{
			"id":        postFile.GetID(),
			"name":      postFile.GetTitle(),
			"url":       postFile.GetURL(),
			"type":      postFile.GetType(),
			"size":      postFile.GetSize(),
			"extension": postFile.GetExtension(),
			"sequence":  postFile.GetSequence(),
		})
	}

	return api.SuccessWithData("Files uploaded successfully", map[string]any{
		"files": uploaded,
	}).ToString()
}

func (c *postMediaComponent) HandleSave(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("post_id is required").ToString()
	}

	var reqData struct {
		PostID string `json:"post_id"`
		Files  []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Sequence int    `json:"sequence"`
		} `json:"files"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	for _, item := range reqData.Files {
		file, err := c.app.GetBlogStore().MediaFindByID(r.Context(), item.ID)
		if err != nil || file == nil {
			continue
		}
		if item.Name != "" {
			file.SetTitle(item.Name)
		}
		file.SetSequence(item.Sequence)
		if err := c.app.GetBlogStore().MediaUpdate(r.Context(), file); err != nil {
			slog.Error("Failed to update post file", "error", err, "file_id", item.ID)
		}
	}

	return api.Success("Media saved successfully").ToString()
}

func (c *postMediaComponent) HandleDelete(w http.ResponseWriter, r *http.Request) string {
	fileID := req.GetStringTrimmed(r, "file_id")
	if fileID == "" {
		return api.Error("file_id is required").ToString()
	}

	file, err := c.app.GetBlogStore().MediaFindByID(r.Context(), fileID)
	if err != nil {
		return api.Error("Failed to find file: " + err.Error()).ToString()
	}
	if file == nil {
		return api.Error("File not found").ToString()
	}

	if err := c.app.GetBlogStore().MediaDeleteByID(r.Context(), fileID); err != nil {
		return api.Error("Failed to delete file: " + err.Error()).ToString()
	}

	return api.Success("File deleted successfully").ToString()
}

func (c *postMediaComponent) HandleAdd(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("post_id is required").ToString()
	}

	mediaURL := req.GetStringTrimmed(r, "media_url")
	if mediaURL == "" {
		return api.Error("media_url is required").ToString()
	}

	mediaFileName := req.GetStringTrimmed(r, "media_file_name")
	mediaType := req.GetStringTrimmed(r, "media_type")

	existingFiles, _ := c.app.GetBlogStore().MediaListByEntityID(r.Context(), postID)
	startSequence := len(existingFiles)

	postFile := blogstore.NewMedia().
		SetEntityID(postID).
		SetTitle(mediaFileName).
		SetURL(mediaURL).
		SetType(mediaType).
		SetSequence(startSequence)

	if err := c.app.GetBlogStore().MediaCreate(r.Context(), postFile); err != nil {
		return api.Error("Failed to save media: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Media added successfully", map[string]any{
		"file": map[string]any{
			"id":       postFile.GetID(),
			"name":     postFile.GetTitle(),
			"url":      postFile.GetURL(),
			"type":     postFile.GetType(),
			"sequence": postFile.GetSequence(),
		},
	}).ToString()
}
