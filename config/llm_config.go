package config

// ============================================================================
// == Anthropic
// ============================================================================

// AnthropicApiUsed is a boolean indicating whether the Anthropic API is used.
var AnthropicApiUsed bool = false

// AnthropicApiKey is the Anthropic API key.
var AnthropicApiKey string

// AnthropicDefaultModel is the default model to use for Anthropic.
var AnthropicDefaultModel string = "claude-3-opus-20240229"

// ============================================================================
// == Google Gemini
// ============================================================================

// GoogleGeminiApiUsed is a boolean indicating whether the Google Gemini API is used.
var GoogleGeminiApiUsed bool = false

// GoogleGeminiApiKey is the Google Gemini API key.
var GoogleGeminiApiKey string

// GoogleGeminiDefaultModel is the default model to use for Google Gemini.
var GoogleGeminiDefaultModel string = "gemini-2.0-flash"

// ============================================================================
// == OpenAI
// ============================================================================

// OpenAiApiUsed is a boolean indicating whether the OpenAI API is used.
var OpenAiApiUsed bool = false

// OpenAiApiKey is the OpenAI API key.
var OpenAiApiKey string

// OpenAiDefaultModel is the default model to use for OpenAI.
var OpenAiDefaultModel string = "gpt-3.5-turbo"

// ============================================================================
// == Vertex AI
// ============================================================================

// VertexAiUsed is a boolean indicating whether the Vertex API is used.
var VertexAiUsed bool = false

// VertexAiProjectID is the Vertex project ID.
var VertexAiProjectID string

// VertexAiRegionID is the Vertex region ID.
var VertexAiRegionID string

// VertexAiModelID is the Vertex model ID.
var VertexAiModelID string

// VertexAiDefaultModel is the default model to use for Vertex.
var VertexAiDefaultModel string = "gemini-2.0-flash"
