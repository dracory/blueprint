package blog_settings

import (
	"net/http"
	"os"
	"strings"

	"github.com/dracory/api"
)

func (controller *blogSettingsController) handleFetchData(r *http.Request) string {
	store := controller.app.GetSettingStore()
	if store == nil {
		return api.Error("Setting store is not configured").ToString()
	}

	value, err := store.Get(r.Context(), SettingKeyBlogTopic, "")
	if err != nil {
		controller.app.GetLogger().Error("Blog settings: failed to load blog topic", "error", err.Error())
		return api.Error("Failed to load blog settings").ToString()
	}

	isEnvOverride := false
	infoMessage := ""
	envTopic := strings.TrimSpace(os.Getenv("BLOG_TOPIC"))
	if envTopic != "" {
		value = envTopic
		isEnvOverride = true
		infoMessage = "The BLOG_TOPIC environment variable is set, so updates are disabled here."
	}

	return api.SuccessWithData("Settings loaded", map[string]any{
		"blog_topic":      value,
		"is_env_override": isEnvOverride,
		"info_message":    infoMessage,
	}).ToString()
}
