package config

import "github.com/dracory/env"

// translationConfig captures i18n settings.
type translationConfig struct {
	defaultLanguage string            // Default language code (e.g., "en", "es", "fr")
	languageList    map[string]string // Mapping of language codes to language names
}

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
