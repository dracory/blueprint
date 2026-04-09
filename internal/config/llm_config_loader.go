package config

// loadLLMConfig loads LLM provider configuration directly into the config.
func loadLLMConfig(cfg ConfigInterface, v *envValidator) {
	anthropicUsed := v.GetBool(KEY_ANTHROPIC_API_USED)
	anthropicKey := v.GetString(KEY_ANTHROPIC_API_KEY)
	anthropicDefaultModel := v.GetString(KEY_ANTHROPIC_API_DEFAULT_MODEL)

	googleGeminiUsed := v.GetBool(KEY_GEMINI_API_USED)
	googleGeminiKey := v.GetString(KEY_GEMINI_API_KEY)
	googleGeminiDefaultModel := v.GetString(KEY_GEMINI_API_DEFAULT_MODEL)

	openAiUsed := v.GetBool(KEY_OPENAI_API_USED)
	openAiKey := v.GetString(KEY_OPENAI_API_KEY)
	openAiDefaultModel := v.GetString(KEY_OPENAI_API_DEFAULT_MODEL)

	openRouterUsed := v.GetBool(KEY_OPENROUTER_API_USED)
	openRouterKey := v.GetString(KEY_OPENROUTER_API_KEY)
	openRouterDefaultModel := v.GetString(KEY_OPENROUTER_API_DEFAULT_MODEL)

	vertexAiUsed := v.GetBool(KEY_VERTEX_AI_API_USED)
	vertexAiModelID := v.GetString(KEY_VERTEX_AI_API_MODEL_ID)
	vertexAiProjectID := v.GetString(KEY_VERTEX_AI_API_PROJECT_ID)
	vertexAiRegionID := v.GetString(KEY_VERTEX_AI_API_REGION_ID)
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
