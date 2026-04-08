package config

import "github.com/dracory/env"

// ============================================================================
// Interface
// ============================================================================

// I18nConfigInterface defines internationalization configuration methods.
type I18nConfigInterface interface {
	SetTranslationLanguageDefault(string)
	GetTranslationLanguageDefault() string

	SetTranslationLanguageList(map[string]string)
	GetTranslationLanguageList() map[string]string
}

// ============================================================================
// Types
// ============================================================================

// translationConfig captures i18n settings.
type translationConfig struct {
	defaultLanguage string            // Default language code (e.g., "en", "es", "fr")
	languageList    map[string]string // Mapping of language codes to language names
}

// ============================================================================
// Loader
// ============================================================================

// loadTranslationConfig loads translation configuration from environment variables.
func loadTranslationConfig() translationConfig {
	defaultLang := env.GetString(KEY_TRANSLATION_LANGUAGE_DEFAULT)
	if defaultLang == "" {
		defaultLang = translationLanguageDefault()
	}

	return translationConfig{
		defaultLanguage: defaultLang,
		languageList:    translationLanguageListDefault(),
	}
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetTranslationLanguageDefault(v string) {
	c.translationLanguageDefault = v
}

func (c *configImplementation) GetTranslationLanguageDefault() string {
	return c.translationLanguageDefault
}

func (c *configImplementation) SetTranslationLanguageList(v map[string]string) {
	c.translationLanguageList = v
}

func (c *configImplementation) GetTranslationLanguageList() map[string]string {
	return c.translationLanguageList
}
