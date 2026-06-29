package aititlegenerator

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
)

func (c *AiTitleGeneratorController) handleSettingsFetchData(r *http.Request) string {
	store := c.app.GetSettingStore()
	if store == nil {
		return api.Error("Setting store is not configured").ToString()
	}

	value, err := store.Get(r.Context(), SETTING_KEY_BLOG_TOPIC, "")
	if err != nil {
		if c.app.GetLogger() != nil {
			c.app.GetLogger().Error("AI title generator settings: failed to load blog topic", "error", err.Error())
		}
		return api.Error("Failed to load settings").ToString()
	}

	infoMessage := ""
	if strings.TrimSpace(value) == "" {
		infoMessage = "Set the Title Generator settings first, then you can generate new titles."
	}

	return api.SuccessWithData("Settings loaded", map[string]any{
		"blog_topic":   strings.TrimSpace(value),
		"info_message": infoMessage,
	}).ToString()
}
