package config

import (
	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// llmConfig captures LLM provider settings.
type llmConfig struct {
	anthropicUsed            bool
	anthropicKey             string
	anthropicDefaultModel    string
	googleGeminiUsed         bool
	googleGeminiKey          string
	googleGeminiDefaultModel string
	openAiUsed               bool
	openAiKey                string
	openAiDefaultModel       string
	openRouterUsed           bool
	openRouterKey            string
	openRouterDefaultModel   string
	vertexAiUsed             bool
	vertexAiModelID          string
	vertexAiProjectID        string
	vertexAiRegionID         string
	vertexAiDefaultModel     string
}

// loadLLMConfig loads LLM provider configuration from environment variables.
func loadLLMConfig(acc *baseCfg.LoadAccumulator) llmConfig {
	anthropicUsed := env.GetBool(KEY_ANTHROPIC_API_USED)
	anthropicKey := env.GetString(KEY_ANTHROPIC_API_KEY)
	anthropicDefaultModel := env.GetString(KEY_ANTHROPIC_API_DEFAULT_MODEL)

	googleGeminiUsed := env.GetBool(KEY_GEMINI_API_USED)
	googleGeminiKey := env.GetString(KEY_GEMINI_API_KEY)
	googleGeminiDefaultModel := env.GetString(KEY_GEMINI_API_DEFAULT_MODEL)

	openAiUsed := env.GetBool(KEY_OPENAI_API_USED)
	openAiKey := env.GetString(KEY_OPENAI_API_KEY)
	openAiDefaultModel := env.GetString(KEY_OPENAI_API_DEFAULT_MODEL)

	openRouterUsed := env.GetBool(KEY_OPENROUTER_API_USED)
	openRouterKey := env.GetString(KEY_OPENROUTER_API_KEY)
	openRouterDefaultModel := env.GetString(KEY_OPENROUTER_API_DEFAULT_MODEL)

	vertexAiUsed := env.GetBool(KEY_VERTEX_AI_API_USED)
	vertexAiModelID := env.GetString(KEY_VERTEX_AI_API_MODEL_ID)
	vertexAiProjectID := env.GetString(KEY_VERTEX_AI_API_PROJECT_ID)
	vertexAiRegionID := env.GetString(KEY_VERTEX_AI_API_REGION_ID)
	vertexAiDefaultModel := env.GetString(KEY_VERTEX_AI_API_DEFAULT_MODEL)

	acc.MustWhen(anthropicUsed, KEY_ANTHROPIC_API_KEY,
		"required when `ANTHROPIC_API_USED` is true", anthropicKey)
	acc.MustWhen(anthropicUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL,
		"required when `ANTHROPIC_API_USED` is true", anthropicDefaultModel)

	acc.MustWhen(googleGeminiUsed, KEY_GEMINI_API_KEY,
		"required when `GEMINI_API_USED` is true", googleGeminiKey)
	acc.MustWhen(googleGeminiUsed, KEY_GEMINI_API_DEFAULT_MODEL,
		"required when `GEMINI_API_USED` is true", googleGeminiDefaultModel)

	acc.MustWhen(openAiUsed, KEY_OPENAI_API_KEY,
		"required when `OPENAI_API_USED` is true", openAiKey)
	acc.MustWhen(openAiUsed, KEY_OPENAI_API_DEFAULT_MODEL,
		"required when `OPENAI_API_USED` is true", openAiDefaultModel)

	acc.MustWhen(openRouterUsed, KEY_OPENROUTER_API_KEY,
		"required when `OPENROUTER_API_USED` is true", openRouterKey)
	acc.MustWhen(openRouterUsed, KEY_OPENROUTER_API_DEFAULT_MODEL,
		"required when `OPENROUTER_API_USED` is true", openRouterDefaultModel)

	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_MODEL_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiModelID)
	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_PROJECT_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiProjectID)
	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_REGION_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiRegionID)
	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL,
		"required when `VERTEX_AI_API_USED` is true", vertexAiDefaultModel)

	return llmConfig{
		anthropicUsed:            anthropicUsed,
		anthropicKey:             anthropicKey,
		anthropicDefaultModel:    anthropicDefaultModel,
		googleGeminiUsed:         googleGeminiUsed,
		googleGeminiKey:          googleGeminiKey,
		googleGeminiDefaultModel: googleGeminiDefaultModel,
		openAiUsed:               openAiUsed,
		openAiKey:                openAiKey,
		openAiDefaultModel:       openAiDefaultModel,
		openRouterUsed:           openRouterUsed,
		openRouterKey:            openRouterKey,
		openRouterDefaultModel:   openRouterDefaultModel,
		vertexAiUsed:             vertexAiUsed,
		vertexAiModelID:          vertexAiModelID,
		vertexAiProjectID:        vertexAiProjectID,
		vertexAiRegionID:         vertexAiRegionID,
		vertexAiDefaultModel:     vertexAiDefaultModel,
	}
}
