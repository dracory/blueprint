package config

// readLLMConfig reads LLM provider configuration from environment variables.
// Each provider is opt-in: set the *_USED flag to true and provide the required keys.
func llmConfig(cfg *configImplementation, v *envValidator) {
	// -------------------------------------------------------------------------
	// Anthropic (Claude)
	// https://console.anthropic.com/settings/keys
	// -------------------------------------------------------------------------

	// Set ANTHROPIC_API_USED=true to enable Anthropic Claude models.
	anthropicUsed := v.GetBool(KEY_ANTHROPIC_API_USED)

	// Your Anthropic API key. Required when ANTHROPIC_API_USED=true.
	anthropicKey := v.GetString(KEY_ANTHROPIC_API_KEY)

	// Default Claude model to use. Required when ANTHROPIC_API_USED=true.
	// Example: claude-3-5-sonnet-20241022
	anthropicDefaultModel := v.GetString(KEY_ANTHROPIC_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// Google Gemini
	// https://aistudio.google.com/app/apikey
	// -------------------------------------------------------------------------

	// Set GEMINI_API_USED=true to enable Google Gemini models.
	googleGeminiUsed := v.GetBool(KEY_GEMINI_API_USED)

	// Your Google Gemini API key. Required when GEMINI_API_USED=true.
	googleGeminiKey := v.GetString(KEY_GEMINI_API_KEY)

	// Default Gemini model to use. Required when GEMINI_API_USED=true.
	// Example: gemini-1.5-pro
	googleGeminiDefaultModel := v.GetString(KEY_GEMINI_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// OpenAI
	// https://platform.openai.com/api-keys
	// -------------------------------------------------------------------------

	// Set OPENAI_API_USED=true to enable OpenAI models.
	openAiUsed := v.GetBool(KEY_OPENAI_API_USED)

	// Your OpenAI API key. Required when OPENAI_API_USED=true.
	openAiKey := v.GetString(KEY_OPENAI_API_KEY)

	// Default OpenAI model to use. Required when OPENAI_API_USED=true.
	// Example: gpt-4o, gpt-4o-mini
	openAiDefaultModel := v.GetString(KEY_OPENAI_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// OpenRouter (multi-provider gateway)
	// https://openrouter.ai/keys
	// -------------------------------------------------------------------------

	// Set OPENROUTER_API_USED=true to enable OpenRouter.
	openRouterUsed := v.GetBool(KEY_OPENROUTER_API_USED)

	// Your OpenRouter API key. Required when OPENROUTER_API_USED=true.
	openRouterKey := v.GetString(KEY_OPENROUTER_API_KEY)

	// Default model to use via OpenRouter. Required when OPENROUTER_API_USED=true.
	// Example: openai/gpt-4o, anthropic/claude-3-5-sonnet
	openRouterDefaultModel := v.GetString(KEY_OPENROUTER_API_DEFAULT_MODEL)

	// -------------------------------------------------------------------------
	// Google Vertex AI
	// https://cloud.google.com/vertex-ai
	// -------------------------------------------------------------------------

	// Set VERTEX_AI_API_USED=true to enable Google Vertex AI.
	vertexAiUsed := v.GetBool(KEY_VERTEX_AI_API_USED)

	// Vertex AI model ID. Required when VERTEX_AI_API_USED=true.
	// Example: gemini-1.5-pro-002
	vertexAiModelID := v.GetString(KEY_VERTEX_AI_API_MODEL_ID)

	// Google Cloud project ID. Required when VERTEX_AI_API_USED=true.
	vertexAiProjectID := v.GetString(KEY_VERTEX_AI_API_PROJECT_ID)

	// Google Cloud region. Required when VERTEX_AI_API_USED=true.
	// Example: us-central1
	vertexAiRegionID := v.GetString(KEY_VERTEX_AI_API_REGION_ID)

	// Default Vertex AI model. Required when VERTEX_AI_API_USED=true.
	vertexAiDefaultModel := v.GetString(KEY_VERTEX_AI_API_DEFAULT_MODEL)

	v.RequireWhen(anthropicUsed, KEY_ANTHROPIC_API_KEY,
		"required when `ANTHROPIC_API_USED` is true", anthropicKey)
	v.RequireWhen(anthropicUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL,
		"required when `ANTHROPIC_API_USED` is true", anthropicDefaultModel)

	v.RequireWhen(googleGeminiUsed, KEY_GEMINI_API_KEY,
		"required when `GEMINI_API_USED` is true", googleGeminiKey)
	v.RequireWhen(googleGeminiUsed, KEY_GEMINI_API_DEFAULT_MODEL,
		"required when `GEMINI_API_USED` is true", googleGeminiDefaultModel)

	v.RequireWhen(openAiUsed, KEY_OPENAI_API_KEY,
		"required when `OPENAI_API_USED` is true", openAiKey)
	v.RequireWhen(openAiUsed, KEY_OPENAI_API_DEFAULT_MODEL,
		"required when `OPENAI_API_USED` is true", openAiDefaultModel)

	v.RequireWhen(openRouterUsed, KEY_OPENROUTER_API_KEY,
		"required when `OPENROUTER_API_USED` is true", openRouterKey)
	v.RequireWhen(openRouterUsed, KEY_OPENROUTER_API_DEFAULT_MODEL,
		"required when `OPENROUTER_API_USED` is true", openRouterDefaultModel)

	v.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_MODEL_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiModelID)
	v.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_PROJECT_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiProjectID)
	v.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_REGION_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiRegionID)
	v.RequireWhen(vertexAiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL,
		"required when `VERTEX_AI_API_USED` is true", vertexAiDefaultModel)

	// -------------------------------------------------------------------------
	// Do not edit below this line
	// -------------------------------------------------------------------------
	cfg.setLLMConfig(
		anthropicUsed, anthropicKey, anthropicDefaultModel,
		googleGeminiUsed, googleGeminiKey, googleGeminiDefaultModel,
		openAiUsed, openAiKey, openAiDefaultModel,
		openRouterUsed, openRouterKey, openRouterDefaultModel,
		vertexAiUsed, vertexAiModelID, vertexAiProjectID, vertexAiRegionID, vertexAiDefaultModel,
	)
}
