package config

// loadLLMConfig loads LLM provider configuration directly into the config.
// Each provider is opt-in: set the *_USED flag to true and provide the required keys.
func loadLLMConfig(cfg ConfigInterface, v *envValidator) {
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

	v.MustWhen(anthropicUsed, KEY_ANTHROPIC_API_KEY,
		"required when `ANTHROPIC_API_USED` is true", anthropicKey)
	v.MustWhen(anthropicUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL,
		"required when `ANTHROPIC_API_USED` is true", anthropicDefaultModel)

	v.MustWhen(googleGeminiUsed, KEY_GEMINI_API_KEY,
		"required when `GEMINI_API_USED` is true", googleGeminiKey)
	v.MustWhen(googleGeminiUsed, KEY_GEMINI_API_DEFAULT_MODEL,
		"required when `GEMINI_API_USED` is true", googleGeminiDefaultModel)

	v.MustWhen(openAiUsed, KEY_OPENAI_API_KEY,
		"required when `OPENAI_API_USED` is true", openAiKey)
	v.MustWhen(openAiUsed, KEY_OPENAI_API_DEFAULT_MODEL,
		"required when `OPENAI_API_USED` is true", openAiDefaultModel)

	v.MustWhen(openRouterUsed, KEY_OPENROUTER_API_KEY,
		"required when `OPENROUTER_API_USED` is true", openRouterKey)
	v.MustWhen(openRouterUsed, KEY_OPENROUTER_API_DEFAULT_MODEL,
		"required when `OPENROUTER_API_USED` is true", openRouterDefaultModel)

	v.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_MODEL_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiModelID)
	v.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_PROJECT_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiProjectID)
	v.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_REGION_ID,
		"required when `VERTEX_AI_API_USED` is true", vertexAiRegionID)
	v.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL,
		"required when `VERTEX_AI_API_USED` is true", vertexAiDefaultModel)

	cfg.SetAnthropicApiUsed(anthropicUsed)
	cfg.SetAnthropicApiKey(anthropicKey)
	cfg.SetAnthropicApiDefaultModel(anthropicDefaultModel)

	cfg.SetGoogleGeminiApiUsed(googleGeminiUsed)
	cfg.SetGoogleGeminiApiKey(googleGeminiKey)
	cfg.SetGoogleGeminiApiDefaultModel(googleGeminiDefaultModel)

	cfg.SetOpenAiApiUsed(openAiUsed)
	cfg.SetOpenAiApiKey(openAiKey)
	cfg.SetOpenAiApiDefaultModel(openAiDefaultModel)

	cfg.SetOpenRouterApiUsed(openRouterUsed)
	cfg.SetOpenRouterApiKey(openRouterKey)
	cfg.SetOpenRouterApiDefaultModel(openRouterDefaultModel)

	cfg.SetVertexAiApiUsed(vertexAiUsed)
	cfg.SetVertexAiApiModelID(vertexAiModelID)
	cfg.SetVertexAiApiProjectID(vertexAiProjectID)
	cfg.SetVertexAiApiRegionID(vertexAiRegionID)
	cfg.SetVertexAiApiDefaultModel(vertexAiDefaultModel)
}
