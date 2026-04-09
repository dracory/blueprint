package config

import "github.com/dracory/env"

// readTranslationConfig reads i18n configuration from environment variables.
func readTranslationConfig(cfg *configImplementation) {
	// Default Language
	//
	// The default locale used when no language is specified.
	// Must match one of the keys in the language list below.
	// Example: en, fr, de, es
	defaultLanguage := env.GetStringOrDefault(KEY_TRANSLATION_LANGUAGE_DEFAULT, translationLanguageDefault())

	// Language List
	//
	// The full list of supported locales mapped to their display names.
	// To add a language, update translationLanguageListDefault() in defaults.go.
	// Example: {"en": "English", "fr": "French"}
	languageList := translationLanguageListDefault()

	cfg.setTranslationConfig(defaultLanguage, languageList)
}
