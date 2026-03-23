package config

import "github.com/dracory/env"

// translationConfig captures i18n settings.
// It manages internationalization configuration including default language
// and available language mappings for the application.
type translationConfig struct {
	defaultLanguage string            // Default language code (e.g., "en", "es", "fr")
	languageList    map[string]string // Mapping of language codes to language names
}

// loadTranslationConfig loads translation configuration from environment variables.
// It reads the default language setting and loads the available language list.
// If no default language is specified, it falls back to the system default.
//
// Returns:
//   - translationConfig: Populated configuration struct with i18n settings
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
