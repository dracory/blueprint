package config

// llmConfig reads LLM provider configuration from environment variables.
// Each provider is opt-in: set the *_USED flag to true and provide the required keys.
func llmConfig(env *envValidator) llmSettings {
	// -------------------------------------------------------------------------
	// Anthropic (Claude)
	// https://console.anthropic.com/settings/keys
	// -------------------------------------------------------------------------

	// Set ANTHROPIC_API_USED=true to enable Anthropic Claude models.
	anthropicUsed := env.GetBool(KEY_ANTHROPIC_API_USED)

	// Your Anthropic API key. Required when ANTHROPIC_API_USED=true.
	anthropicKey := env.GetString(KEY_ANTHROPIC_API_KEY)

	// Default Claude model to use. Required when ANTHROPIC_API_USED=true.
	// Example: claude-3-5-sonnet-20241022
	anthropicDefaultModel := env.GetString(KEY_ANTHROPIC_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// Google Gemini
	// https://aistudio.google.com/app/apikey
	// -------------------------------------------------------------------------

	// Set GEMINI_API_USED=true to enable Google Gemini models.
	googleGeminiUsed := env.GetBool(KEY_GEMINI_API_USED)

	// Your Google Gemini API key. Required when GEMINI_API_USED=true.
	googleGeminiKey := env.GetString(KEY_GEMINI_API_KEY)

	// Default Gemini model to use. Required when GEMINI_API_USED=true.
	// Example: gemini-1.5-pro
	googleGeminiDefaultModel := env.GetString(KEY_GEMINI_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// OpenAI
	// https://platform.openai.com/api-keys
	// -------------------------------------------------------------------------

	// Set OPENAI_API_USED=true to enable OpenAI models.
	openAiUsed := env.GetBool(KEY_OPENAI_API_USED)

	// Your OpenAI API key. Required when OPENAI_API_USED=true.
	openAiKey := env.GetString(KEY_OPENAI_API_KEY)

	// Default OpenAI model to use. Required when OPENAI_API_USED=true.
	// Example: gpt-4o, gpt-4o-mini
	openAiDefaultModel := env.GetString(KEY_OPENAI_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// OpenRouter (multi-provider gateway)
	// https://openrouter.ai/keys
	// -------------------------------------------------------------------------

	// Set OPENROUTER_API_USED=true to enable OpenRouter.
	openRouterUsed := env.GetBool(KEY_OPENROUTER_API_USED)

	// Your OpenRouter API key. Required when OPENROUTER_API_USED=true.
	openRouterKey := env.GetString(KEY_OPENROUTER_API_KEY)

	// Default model to use via OpenRouter. Required when OPENROUTER_API_USED=true.
	// Example: openai/gpt-4o, anthropic/claude-3-5-sonnet
	openRouterDefaultModel := env.GetString(KEY_OPENROUTER_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// Google Vertex AI
	// https://cloud.google.com/vertex-ai
	// -------------------------------------------------------------------------

	// Set VERTEX_AI_API_USED=true to enable Google Vertex AI.
	vertexAiUsed := env.GetBool(KEY_VERTEX_AI_API_USED)

	// Vertex AI model ID. Required when VERTEX_AI_API_USED=true.
	// Example: gemini-1.5-pro-002
	vertexAiModelID := env.GetString(KEY_VERTEX_AI_API_MODEL_ID)

	// Google Cloud project ID. Required when VERTEX_AI_API_USED=true.
	vertexAiProjectID := env.GetString(KEY_VERTEX_AI_API_PROJECT_ID)

	// Google Cloud region. Required when VERTEX_AI_API_USED=true.
	// Example: us-central1
	vertexAiRegionID := env.GetString(KEY_VERTEX_AI_API_REGION_ID)

	// Default Vertex AI model. Required when VERTEX_AI_API_USED=true.
	vertexAiDefaultModel := env.GetString(KEY_VERTEX_AI_API_DEFAULT_MODEL)

	env.RequireWhen(anthropicUsed, KEY_ANTHROPIC_API_KEY,
		"required when `ANTHROPIC_API_USED` is true", anthropicKey)
	env.RequireWhen(anthropicUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL,
		"required when `ANTHROPIC_API_USED` is true", anthropicDefaultModel)

	env.RequireWhen(googleGeminiUsed, KEY_GEMINI_API_KEY,
		"required when `GEMINI_API_USED` is true", googleGeminiKey)
	env.RequireWhen(googleGeminiUsed, KEY_GEMINI_API_DEFAULT_MODEL,
		"required when `GEMINI_API_USED` is true", googleGeminiDefaultModel)

	env.RequireWhen(openAiUsed, KEY_OPENAI_API_KEY,
		"required when `OPENAI_API_USED` is true", openAiKey)
	env.RequireWhen(openAiUsed, KEY_OPENAI_API_DEFAULT_MODEL,
		"required when `OPENAI_API_USED` is true", openAiDefaultModel)

	env.RequireWhen(openRouterUsed, KEY_OPENROUTER_API_KEY,
		"required when `OPENROUTER_API_USED` is true", openRouterKey)
	env.RequireWhen(openRouterUsed, KEY_OPENROUTER_API_DEFAULT_MODEL,
		"required when `OPENROUTER_API_USED` is true", openRouterDefaultModel)

	env.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_MODEL_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiModelID)
	env.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_PROJECT_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiProjectID)
	env.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_REGION_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiRegionID)
	env.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL,
		"required when `VERTEX_AI_API_USED` is true", vertexAiDefaultModel)

	return llmSettings{
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

type llmSettings struct {
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
