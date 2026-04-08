package config

import (
	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// ============================================================================
// Interface
// ============================================================================

// LLMConfigInterface defines LLM provider configuration methods.
type LLMConfigInterface interface {
	// Anthropic
	SetAnthropicApiUsed(bool)
	GetAnthropicApiUsed() bool

	SetAnthropicApiKey(string)
	GetAnthropicApiKey() string

	SetAnthropicApiDefaultModel(string)
	GetAnthropicApiDefaultModel() string

	// Google Gemini
	SetGoogleGeminiApiUsed(bool)
	GetGoogleGeminiApiUsed() bool

	SetGoogleGeminiApiKey(string)
	GetGoogleGeminiApiKey() string

	SetGoogleGeminiApiDefaultModel(string)
	GetGoogleGeminiApiDefaultModel() string

	// OpenRouter
	SetOpenRouterApiUsed(bool)
	GetOpenRouterApiUsed() bool

	SetOpenRouterApiKey(string)
	GetOpenRouterApiKey() string

	SetOpenRouterApiDefaultModel(string)
	GetOpenRouterApiDefaultModel() string

	// OpenAI
	SetOpenAiApiUsed(bool)
	GetOpenAiApiUsed() bool

	SetOpenAiApiKey(string)
	GetOpenAiApiKey() string

	SetOpenAiApiDefaultModel(string)
	GetOpenAiApiDefaultModel() string

	// Vertex AI
	SetVertexAiApiUsed(bool)
	GetVertexAiApiUsed() bool

	SetVertexAiApiDefaultModel(string)
	GetVertexAiApiDefaultModel() string

	SetVertexAiApiProjectID(string)
	GetVertexAiApiProjectID() string

	SetVertexAiApiRegionID(string)
	GetVertexAiApiRegionID() string

	SetVertexAiApiModelID(string)
	GetVertexAiApiModelID() string
}

// ============================================================================
// Types
// ============================================================================

// llmConfig captures LLM provider settings.
type llmConfig struct {
	anthropicUsed         bool   // Enable Anthropic Claude API
	anthropicKey          string // Anthropic API key
	anthropicDefaultModel string // Default Anthropic model (e.g., claude-3-sonnet)

	googleGeminiUsed         bool   // Enable Google Gemini API
	googleGeminiKey          string // Google Gemini API key
	googleGeminiDefaultModel string // Default Gemini model (e.g., gemini-pro)

	openAiUsed         bool   // Enable OpenAI API
	openAiKey          string // OpenAI API key
	openAiDefaultModel string // Default OpenAI model (e.g., gpt-4)

	openRouterUsed         bool   // Enable OpenRouter API
	openRouterKey          string // OpenRouter API key
	openRouterDefaultModel string // Default OpenRouter model

	vertexAiUsed         bool   // Enable Google Vertex AI API
	vertexAiModelID      string // Vertex AI model identifier
	vertexAiProjectID    string // Google Cloud project ID
	vertexAiRegionID     string // Google Cloud region ID
	vertexAiDefaultModel string // Default Vertex AI model
}

// ============================================================================
// Loader
// ============================================================================

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

	acc.MustWhen(anthropicUsed, KEY_ANTHROPIC_API_KEY, "required when `ANTHROPIC_API_USED` is true", anthropicKey)
	acc.MustWhen(anthropicUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL, "required when `ANTHROPIC_API_USED` is true", anthropicDefaultModel)

	acc.MustWhen(googleGeminiUsed, KEY_GEMINI_API_KEY, "required when `GEMINI_API_USED` is true", googleGeminiKey)
	acc.MustWhen(googleGeminiUsed, KEY_GEMINI_API_DEFAULT_MODEL, "required when `GEMINI_API_USED` is true", googleGeminiDefaultModel)

	acc.MustWhen(openAiUsed, KEY_OPENAI_API_KEY, "required when `OPENAI_API_USED` is true", openAiKey)
	acc.MustWhen(openAiUsed, KEY_OPENAI_API_DEFAULT_MODEL, "required when `OPENAI_API_USED` is true", openAiDefaultModel)

	acc.MustWhen(openRouterUsed, KEY_OPENROUTER_API_KEY, "required when `OPENROUTER_API_USED` is true", openRouterKey)
	acc.MustWhen(openRouterUsed, KEY_OPENROUTER_API_DEFAULT_MODEL, "required when `OPENROUTER_API_USED` is true", openRouterDefaultModel)

	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_MODEL_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiModelID)
	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_PROJECT_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiProjectID)
	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_REGION_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiRegionID)
	acc.MustWhen(vertexAiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL, "required when `VERTEX_AI_API_USED` is true", vertexAiDefaultModel)

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

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

// Anthropic
func (c *configImplementation) SetAnthropicApiUsed(v bool) {
	c.anthropicApiUsed = v
}

func (c *configImplementation) GetAnthropicApiUsed() bool {
	return c.anthropicApiUsed
}

func (c *configImplementation) SetAnthropicApiKey(v string) {
	c.anthropicApiKey = v
}

func (c *configImplementation) GetAnthropicApiKey() string {
	return c.anthropicApiKey
}

func (c *configImplementation) SetAnthropicApiDefaultModel(v string) {
	c.anthropicApiDefaultModel = v
}

func (c *configImplementation) GetAnthropicApiDefaultModel() string {
	return c.anthropicApiDefaultModel
}

// Google Gemini
func (c *configImplementation) SetGoogleGeminiApiUsed(v bool) {
	c.googleGeminiApiUsed = v
}

func (c *configImplementation) GetGoogleGeminiApiUsed() bool {
	return c.googleGeminiApiUsed
}

func (c *configImplementation) SetGoogleGeminiApiKey(v string) {
	c.googleGeminiApiKey = v
}

func (c *configImplementation) GetGoogleGeminiApiKey() string {
	return c.googleGeminiApiKey
}

func (c *configImplementation) SetGoogleGeminiApiDefaultModel(v string) {
	c.googleGeminiApiDefaultModel = v
}

func (c *configImplementation) GetGoogleGeminiApiDefaultModel() string {
	return c.googleGeminiApiDefaultModel
}

// OpenAI
func (c *configImplementation) SetOpenAiApiUsed(v bool) {
	c.openAiApiUsed = v
}

func (c *configImplementation) GetOpenAiApiUsed() bool {
	return c.openAiApiUsed
}

func (c *configImplementation) SetOpenAiApiKey(v string) {
	c.openAiApiKey = v
}

func (c *configImplementation) GetOpenAiApiKey() string {
	return c.openAiApiKey
}

func (c *configImplementation) SetOpenAiApiDefaultModel(v string) {
	c.openAiApiDefaultModel = v
}

func (c *configImplementation) GetOpenAiApiDefaultModel() string {
	return c.openAiApiDefaultModel
}

// OpenRouter
func (c *configImplementation) SetOpenRouterApiKey(v string) {
	c.openRouterApiKey = v
}

func (c *configImplementation) GetOpenRouterApiKey() string {
	return c.openRouterApiKey
}

func (c *configImplementation) SetOpenRouterApiUsed(v bool) {
	c.openRouterApiUsed = v
}

func (c *configImplementation) GetOpenRouterApiUsed() bool {
	return c.openRouterApiUsed
}

func (c *configImplementation) SetOpenRouterApiDefaultModel(v string) {
	c.openRouterApiDefaultModel = v
}

func (c *configImplementation) GetOpenRouterApiDefaultModel() string {
	return c.openRouterApiDefaultModel
}

// Vertex AI
func (c *configImplementation) SetVertexAiApiUsed(v bool) {
	c.vertexAiApiUsed = v
}

func (c *configImplementation) GetVertexAiApiUsed() bool {
	return c.vertexAiApiUsed
}

func (c *configImplementation) SetVertexAiApiDefaultModel(v string) {
	c.vertexAiApiDefaultModel = v
}

func (c *configImplementation) GetVertexAiApiDefaultModel() string {
	return c.vertexAiApiDefaultModel
}

func (c *configImplementation) SetVertexAiApiProjectID(v string) {
	c.vertexAiApiProjectID = v
}

func (c *configImplementation) GetVertexAiApiProjectID() string {
	return c.vertexAiApiProjectID
}

func (c *configImplementation) SetVertexAiApiRegionID(v string) {
	c.vertexAiApiRegionID = v
}

func (c *configImplementation) GetVertexAiApiRegionID() string {
	return c.vertexAiApiRegionID
}

func (c *configImplementation) SetVertexAiApiModelID(v string) {
	c.vertexAiApiModelID = v
}

func (c *configImplementation) GetVertexAiApiModelID() string {
	return c.vertexAiApiModelID
}
