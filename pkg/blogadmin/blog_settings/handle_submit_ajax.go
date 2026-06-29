package blog_settings

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"project/pkg/blogadmin/shared"

	"github.com/dracory/api"
)

func (controller *blogSettingsController) handleSubmit(r *http.Request) string {
	var reqBody struct {
		Action    string `json:"action"`
		BlogTopic string `json:"blog_topic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	topic := strings.TrimSpace(reqBody.BlogTopic)
	if topic == "" {
		return api.Error("Blog topic is required").ToString()
	}

	if strings.TrimSpace(os.Getenv("BLOG_TOPIC")) != "" {
		return api.Error("Blog topic is managed via environment and cannot be changed here.").ToString()
	}

	store := controller.app.GetSettingStore()
	if store == nil {
		return api.Error("Setting store is not configured").ToString()
	}

	if err := store.Set(r.Context(), SettingKeyBlogTopic, topic); err != nil {
		controller.app.GetLogger().Error("Blog settings: failed to save blog topic", "error", err.Error())
		return api.Error("Failed to save blog topic. Please try again later.").ToString()
	}

	if reqBody.Action == "save_close" {
		return api.SuccessWithData("Blog settings saved successfully", map[string]any{
			"redirect_url": shared.NewLinks("/admin/blog").PostManager(),
		}).ToString()
	}

	return api.Success("Blog settings saved successfully").ToString()
}
